package kubecontroller

import (
	"context"
	"fmt"
	"time"

	brucov1alpha1 "github.com/ferama/bruco/pkg/kube/apis/brucocontroller/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	clientset "github.com/ferama/bruco/pkg/kube/generated/clientset/versioned"
	brucoscheme "github.com/ferama/bruco/pkg/kube/generated/clientset/versioned/scheme"
	informers "github.com/ferama/bruco/pkg/kube/generated/informers/externalversions/brucocontroller/v1alpha1"
	listers "github.com/ferama/bruco/pkg/kube/generated/listers/brucocontroller/v1alpha1"
)

// BrucoController is the controller implementation for Bruco resources
type BrucoProjectController struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// brucoclientset is a clientset for our own API group
	brucoclientset clientset.Interface

	brucosProjectLister listers.BrucoProjectLister
	brucosProjectSynced cache.InformerSynced

	brucoLister  listers.BrucoLister
	brucosSynced cache.InformerSynced

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
func NewBrucoProjectController(
	kubeclientset kubernetes.Interface,
	brucoclientset clientset.Interface,
	brucoInformer informers.BrucoInformer,
	brucoProjectInformer informers.BrucoProjectInformer) *BrucoProjectController {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(brucoscheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &BrucoProjectController{
		kubeclientset:  kubeclientset,
		brucoclientset: brucoclientset,

		brucosProjectLister: brucoProjectInformer.Lister(),
		brucosProjectSynced: brucoProjectInformer.Informer().HasSynced,

		brucoLister:  brucoInformer.Lister(),
		brucosSynced: brucoInformer.Informer().HasSynced,

		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Brucos"),
		recorder:  recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Bruco resources change
	brucoProjectInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueBrucoProject,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueBrucoProject(new)
		},
	})

	brucoInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newBruco := new.(*brucov1alpha1.Bruco)
			oldBruco := old.(*brucov1alpha1.Bruco)
			if newBruco.ResourceVersion == oldBruco.ResourceVersion {
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
func (c *BrucoProjectController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting BrucoProject controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.brucosProjectSynced); !ok {
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
func (c *BrucoProjectController) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *BrucoProjectController) processNextWorkItem() bool {
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
func (c *BrucoProjectController) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Bruco resource with this namespace/name
	brucoProject, err := c.brucosProjectLister.BrucoProjects(namespace).Get(name)
	if err != nil {
		// The Bruco resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {

			utilruntime.HandleError(fmt.Errorf("brucoproject '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	// Create required but not existsing brucos
	for i, brucoConf := range brucoProject.Spec.Brucos {
		brucoName := fmt.Sprintf("%s-%d", brucoProject.Name, i)
		// brucoName := brucoProject.Name
		bruco, err := c.brucoLister.Brucos(brucoProject.Namespace).Get(brucoName)
		if errors.IsNotFound(err) {
			bruco, err = c.brucoclientset.
				BrucoV1alpha1().
				Brucos(brucoProject.Namespace).
				Create(context.TODO(), newBrucoFromProject(
					brucoProject,
					brucoConf,
					brucoName), metav1.CreateOptions{})
		}
		if err != nil {
			return err
		}
		if !metav1.IsControlledBy(bruco, brucoProject) {
			msg := fmt.Sprintf(MessageResourceExists, bruco.Name)
			c.recorder.Event(bruco, corev1.EventTypeWarning, ErrResourceExists, msg)
			return fmt.Errorf(msg)
		}
	}

	c.recorder.Event(brucoProject, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// enqueueBrucoProject takes a Bruco resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Bruco.
func (c *BrucoProjectController) enqueueBrucoProject(obj interface{}) {
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
func (c *BrucoProjectController) handleObject(obj interface{}) {
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
		if ownerRef.Kind != "BrucoProject" {
			return
		}

		brucoProject, err := c.brucosProjectLister.
			BrucoProjects(object.GetNamespace()).
			Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of bruco project '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueBrucoProject(brucoProject)
		return
	}
}
