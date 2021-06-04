package kubecontroller

import (
	"reflect"
	"testing"
	"time"

	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/diff"
	kubeinformers "k8s.io/client-go/informers"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"

	brucocontroller "github.com/ferama/bruco/pkg/kube/apis/brucocontroller/v1alpha1"
	"github.com/ferama/bruco/pkg/kube/generated/clientset/versioned/fake"
	informers "github.com/ferama/bruco/pkg/kube/generated/informers/externalversions"
)

var (
	alwaysReady        = func() bool { return true }
	noResyncPeriodFunc = func() time.Duration { return 0 }
)

type fixture struct {
	t *testing.T

	client     *fake.Clientset
	kubeclient *k8sfake.Clientset
	// Objects to put in the store.
	brucoLister      []*brucocontroller.Bruco
	deploymentLister []*apps.Deployment
	serviceLister    []*corev1.Service
	// Actions expected to happen on the client.
	kubeactions []core.Action
	actions     []core.Action
	// Objects from here preloaded into NewSimpleFake.
	kubeobjects []runtime.Object
	objects     []runtime.Object
}

func newFixture(t *testing.T) *fixture {
	f := &fixture{}
	f.t = t
	f.objects = []runtime.Object{}
	f.kubeobjects = []runtime.Object{}
	return f
}

func newBruco(name string, replicas *int32) *brucocontroller.Bruco {
	return &brucocontroller.Bruco{
		TypeMeta: metav1.TypeMeta{APIVersion: brucocontroller.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: brucocontroller.BrucoSpec{
			// DeploymentName: fmt.Sprintf("%s-deployment", name),
			Replicas: replicas,
		},
	}
}

func (f *fixture) newController() (*Controller, informers.SharedInformerFactory, kubeinformers.SharedInformerFactory) {
	f.client = fake.NewSimpleClientset(f.objects...)
	f.kubeclient = k8sfake.NewSimpleClientset(f.kubeobjects...)

	i := informers.NewSharedInformerFactory(f.client, noResyncPeriodFunc())
	k8sI := kubeinformers.NewSharedInformerFactory(f.kubeclient, noResyncPeriodFunc())

	c := NewController(f.kubeclient, f.client,
		k8sI.Apps().V1().Deployments(),
		k8sI.Core().V1().Services(),
		i.Brucocontroller().V1alpha1().Brucos())

	c.brucosSynced = alwaysReady
	c.deploymentsSynced = alwaysReady
	c.recorder = &record.FakeRecorder{}

	for _, f := range f.brucoLister {
		i.Brucocontroller().V1alpha1().Brucos().Informer().GetIndexer().Add(f)
	}

	for _, d := range f.deploymentLister {
		k8sI.Apps().V1().Deployments().Informer().GetIndexer().Add(d)
	}

	for _, s := range f.serviceLister {
		k8sI.Core().V1().Services().Informer().GetIndexer().Add(s)
	}

	return c, i, k8sI
}

func (f *fixture) run(brucoName string) {
	f.runController(brucoName, true, false)
}

func (f *fixture) runExpectError(brucoName string) {
	f.runController(brucoName, true, true)
}

func (f *fixture) runController(brucoName string, startInformers bool, expectError bool) {
	c, i, k8sI := f.newController()
	if startInformers {
		stopCh := make(chan struct{})
		defer close(stopCh)
		i.Start(stopCh)
		k8sI.Start(stopCh)
	}

	err := c.syncHandler(brucoName)
	if !expectError && err != nil {
		f.t.Errorf("error syncing bruco: %v", err)
	} else if expectError && err == nil {
		f.t.Error("expected error syncing bruco, got nil")
	}

	actions := filterInformerActions(f.client.Actions())
	for i, action := range actions {
		if len(f.actions) < i+1 {
			f.t.Errorf("%d unexpected actions: %+v", len(actions)-len(f.actions), actions[i:])
			break
		}

		expectedAction := f.actions[i]
		checkAction(expectedAction, action, f.t)
	}

	if len(f.actions) > len(actions) {
		f.t.Errorf("%d additional expected actions:%+v", len(f.actions)-len(actions), f.actions[len(actions):])
	}

	k8sActions := filterInformerActions(f.kubeclient.Actions())
	for i, action := range k8sActions {
		if len(f.kubeactions) < i+1 {
			f.t.Errorf("%d unexpected actions: %+v", len(k8sActions)-len(f.kubeactions), k8sActions[i:])
			break
		}

		expectedAction := f.kubeactions[i]
		checkAction(expectedAction, action, f.t)
	}

	if len(f.kubeactions) > len(k8sActions) {
		f.t.Errorf("%d additional expected actions:%+v", len(f.kubeactions)-len(k8sActions), f.kubeactions[len(k8sActions):])
	}
}

// checkAction verifies that expected and actual actions are equal and both have
// same attached resources
func checkAction(expected, actual core.Action, t *testing.T) {
	if !(expected.Matches(actual.GetVerb(), actual.GetResource().Resource) && actual.GetSubresource() == expected.GetSubresource()) {
		t.Errorf("Expected\n\t%#v\ngot\n\t%#v", expected, actual)
		return
	}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Action has wrong type. Expected: %t. Got: %t", expected, actual)
		return
	}

	switch a := actual.(type) {
	case core.CreateActionImpl:
		e, _ := expected.(core.CreateActionImpl)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expObject, object))
		}
	case core.UpdateActionImpl:
		e, _ := expected.(core.UpdateActionImpl)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expObject, object))
		}
	case core.PatchActionImpl:
		e, _ := expected.(core.PatchActionImpl)
		expPatch := e.GetPatch()
		patch := a.GetPatch()

		if !reflect.DeepEqual(expPatch, patch) {
			t.Errorf("Action %s %s has wrong patch\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expPatch, patch))
		}
	default:
		t.Errorf("Uncaptured Action %s %s, you should explicitly add a case to capture it",
			actual.GetVerb(), actual.GetResource().Resource)
	}
}

// filterInformerActions filters list and watch actions for testing resources.
// Since list and watch don't change resource state we can filter it to lower
// nose level in our tests.
func filterInformerActions(actions []core.Action) []core.Action {
	ret := []core.Action{}
	for _, action := range actions {
		if len(action.GetNamespace()) == 0 &&
			(action.Matches("list", "brucos") ||
				action.Matches("watch", "brucos") ||
				action.Matches("list", "deployments") ||
				action.Matches("watch", "deployments") ||
				action.Matches("list", "services") ||
				action.Matches("watch", "services")) {
			continue
		}
		ret = append(ret, action)
	}

	return ret
}

func (f *fixture) expectCreateDeploymentAction(d *apps.Deployment) {
	f.kubeactions = append(f.kubeactions, core.NewCreateAction(schema.GroupVersionResource{Resource: "deployments"}, d.Namespace, d))
}

func (f *fixture) expectCreateServiceAction(s *corev1.Service) {
	f.kubeactions = append(f.kubeactions, core.NewCreateAction(schema.GroupVersionResource{Resource: "services"}, s.Namespace, s))
}

func (f *fixture) expectUpdateDeploymentAction(d *apps.Deployment) {
	f.kubeactions = append(f.kubeactions, core.NewUpdateAction(schema.GroupVersionResource{Resource: "deployments"}, d.Namespace, d))
}

func (f *fixture) expectUpdateBrucoStatusAction(bruco *brucocontroller.Bruco) {
	action := core.NewUpdateAction(schema.GroupVersionResource{Resource: "brucos"}, bruco.Namespace, bruco)
	action.Subresource = "status"
	f.actions = append(f.actions, action)
}

func getKey(bruco *brucocontroller.Bruco, t *testing.T) string {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(bruco)
	if err != nil {
		t.Errorf("Unexpected error getting key for bruco %v: %v", bruco.Name, err)
		return ""
	}
	return key
}

func TestCreatesDeployment(t *testing.T) {
	f := newFixture(t)
	bruco := newBruco("test", int32Ptr(1))

	f.brucoLister = append(f.brucoLister, bruco)
	f.objects = append(f.objects, bruco)

	expDeployment := newDeployment(bruco)
	f.expectCreateDeploymentAction(expDeployment)
	f.expectUpdateBrucoStatusAction(bruco)

	expService := newService(bruco)
	f.expectCreateServiceAction(expService)

	f.run(getKey(bruco, t))
}

func TestDoNothing(t *testing.T) {
	f := newFixture(t)
	bruco := newBruco("test", int32Ptr(1))
	d := newDeployment(bruco)
	s := newService(bruco)

	f.brucoLister = append(f.brucoLister, bruco)
	f.objects = append(f.objects, bruco)
	f.deploymentLister = append(f.deploymentLister, d)
	f.serviceLister = append(f.serviceLister, s)
	f.kubeobjects = append(f.kubeobjects, d)
	f.kubeobjects = append(f.kubeobjects, s)

	f.expectUpdateBrucoStatusAction(bruco)
	f.run(getKey(bruco, t))
}

func TestUpdateBruco(t *testing.T) {
	f := newFixture(t)
	bruco := newBruco("test", int32Ptr(1))
	d := newDeployment(bruco)
	s := newService(bruco)

	// Update replicas
	bruco.Spec.Replicas = int32Ptr(2)
	expDeployment := newDeployment(bruco)

	f.brucoLister = append(f.brucoLister, bruco)
	f.objects = append(f.objects, bruco)
	f.deploymentLister = append(f.deploymentLister, d)
	f.serviceLister = append(f.serviceLister, s)
	f.kubeobjects = append(f.kubeobjects, d)

	f.expectUpdateBrucoStatusAction(bruco)
	f.expectUpdateDeploymentAction(expDeployment)
	f.run(getKey(bruco, t))
}

func TestNotControlledByUs(t *testing.T) {
	f := newFixture(t)
	bruco := newBruco("test", int32Ptr(1))
	d := newDeployment(bruco)

	d.ObjectMeta.OwnerReferences = []metav1.OwnerReference{}

	f.brucoLister = append(f.brucoLister, bruco)
	f.objects = append(f.objects, bruco)
	f.deploymentLister = append(f.deploymentLister, d)
	f.kubeobjects = append(f.kubeobjects, d)

	f.runExpectError(getKey(bruco, t))
}

func int32Ptr(i int32) *int32 { return &i }
