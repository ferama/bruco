package kubecontroller

import (
	"context"
	"fmt"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	brucov1alpha1 "github.com/ferama/bruco/pkg/kube/apis/brucocontroller/v1alpha1"
	clientset "github.com/ferama/bruco/pkg/kube/generated/clientset/versioned"
	brucoscheme "github.com/ferama/bruco/pkg/kube/generated/clientset/versioned/scheme"
	informers "github.com/ferama/bruco/pkg/kube/generated/informers/externalversions/brucocontroller/v1alpha1"
	listers "github.com/ferama/bruco/pkg/kube/generated/listers/brucocontroller/v1alpha1"
)

// BrucoController is the controller implementation for Bruco resources
type BrucoController struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// brucoclientset is a clientset for our own API group
	brucoclientset clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	servicesLister    corelisters.ServiceLister
	servicesSynced    cache.InformerSynced
	brucosLister      listers.BrucoLister
	brucosSynced      cache.InformerSynced
	configMapLister   corelisters.ConfigMapLister
	configMapSynced   cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewBrucoController returns a new sample controller
func NewBrucoController(
	kubeclientset kubernetes.Interface,
	brucoclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	serviceInformer coreinformers.ServiceInformer,
	configMapInformer coreinformers.ConfigMapInformer,
	brucoInformer informers.BrucoInformer) *BrucoController {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(brucoscheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &BrucoController{
		kubeclientset:     kubeclientset,
		brucoclientset:    brucoclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		servicesLister:    serviceInformer.Lister(),
		servicesSynced:    serviceInformer.Informer().HasSynced,
		configMapLister:   configMapInformer.Lister(),
		configMapSynced:   configMapInformer.Informer().HasSynced,
		brucosLister:      brucoInformer.Lister(),
		brucosSynced:      brucoInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Brucos"),
		recorder:          recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Bruco resources change
	brucoInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueBruco,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueBruco(new)
		},
	})
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a Bruco resource will enqueue that Bruco resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newSvc := new.(*corev1.Service)
			oldSvc := old.(*corev1.Service)
			if newSvc.ResourceVersion == oldSvc.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	configMapInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newCm := new.(*corev1.ConfigMap)
			oldCm := old.(*corev1.ConfigMap)
			if newCm.ResourceVersion == oldCm.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *BrucoController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Bruco controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.brucosSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Bruco resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *BrucoController) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *BrucoController) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Bruco resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Bruco resource
// with the current status of the resource.
func (c *BrucoController) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Bruco resource with this namespace/name
	bruco, err := c.brucosLister.Brucos(namespace).Get(name)
	if err != nil {
		// The Bruco resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("bruco '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	deploymentName := getDeploymentName(bruco)
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}

	// Get the deployment with the name specified in Bruco.spec
	deployment, err := c.deploymentsLister.Deployments(bruco.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.
			AppsV1().
			Deployments(bruco.Namespace).
			Create(context.TODO(), newDeployment(bruco), metav1.CreateOptions{})
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}
	// If the Deployment is not controlled by this Bruco resource, we should log
	// a warning to the event recorder and return error msg.
	if !metav1.IsControlledBy(deployment, bruco) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(bruco, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	serviceName := getServiceName(bruco)
	// Get the service with the name specified in Bruco.spec
	service, err := c.servicesLister.Services(bruco.Namespace).Get(serviceName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		service, err = c.kubeclientset.
			CoreV1().
			Services(bruco.Namespace).
			Create(context.TODO(), newService(bruco), metav1.CreateOptions{})
	}
	if err != nil {
		log.Println(err)
		return err
	}

	if !metav1.IsControlledBy(service, bruco) {
		msg := fmt.Sprintf(MessageResourceExists, service.Name)
		c.recorder.Event(bruco, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	configMapName := getConfigMapName(bruco)
	// Get the service with the name specified in Bruco.spec
	configMap, err := c.configMapLister.ConfigMaps(bruco.Namespace).Get(configMapName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		configMap, err = c.kubeclientset.
			CoreV1().
			ConfigMaps(bruco.Namespace).
			Create(context.TODO(), newConfigMap(bruco), metav1.CreateOptions{})
	}
	if err != nil {
		log.Println(err)
		return err
	}

	if !metav1.IsControlledBy(configMap, bruco) {
		msg := fmt.Sprintf(MessageResourceExists, configMap.Name)
		c.recorder.Event(bruco, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If this number of the replicas on the Bruco resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if bruco.Spec.Replicas != nil && *bruco.Spec.Replicas != *deployment.Spec.Replicas {
		klog.V(4).Infof("Bruco %s replicas: %d, deployment replicas: %d", name, *bruco.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeclientset.AppsV1().
			Deployments(bruco.Namespace).
			Update(context.TODO(), newDeployment(bruco), metav1.UpdateOptions{})

	}

	// restarts deployment on new generation
	if bruco.Generation != bruco.Status.CurrentGeneration {
		_, err = c.kubeclientset.
			CoreV1().
			ConfigMaps(bruco.Namespace).
			Update(context.TODO(), newConfigMap(bruco), metav1.UpdateOptions{})

		if err == nil {
			deployment = newDeployment(bruco)
			if deployment.Annotations == nil {
				deployment.Annotations = make(map[string]string)
			}
			deployment.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
			deployment, err = c.kubeclientset.AppsV1().
				Deployments(bruco.Namespace).
				Update(context.TODO(), deployment, metav1.UpdateOptions{})
		}
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the Bruco resource to reflect the
	// current state of the world
	err = c.updateBrucoStatus(bruco, deployment)
	if err != nil {
		return err
	}

	c.recorder.Event(bruco, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *BrucoController) updateBrucoStatus(bruco *brucov1alpha1.Bruco, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	brucoCopy := bruco.DeepCopy()
	brucoCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	brucoCopy.Status.CurrentGeneration = bruco.Generation
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Bruco resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.brucoclientset.
		BrucoV1alpha1().
		Brucos(bruco.Namespace).
		UpdateStatus(context.TODO(), brucoCopy, metav1.UpdateOptions{})
	return err
}

// enqueueBruco takes a Bruco resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Bruco.
func (c *BrucoController) enqueueBruco(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Bruco resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Bruco resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *BrucoController) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Bruco, we should not do anything more
		// with it.
		if ownerRef.Kind != "Bruco" {
			return
		}

		bruco, err := c.brucosLister.
			Brucos(object.GetNamespace()).
			Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of bruco '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueBruco(bruco)
		return
	}
}
