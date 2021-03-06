// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/ferama/bruco/pkg/kube/apis/brucocontroller/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBrucoProjects implements BrucoProjectInterface
type FakeBrucoProjects struct {
	Fake *FakeBrucoV1alpha1
	ns   string
}

var brucoprojectsResource = schema.GroupVersionResource{Group: "bruco.ferama.github.io", Version: "v1alpha1", Resource: "brucoprojects"}

var brucoprojectsKind = schema.GroupVersionKind{Group: "bruco.ferama.github.io", Version: "v1alpha1", Kind: "BrucoProject"}

// Get takes name of the brucoProject, and returns the corresponding brucoProject object, and an error if there is any.
func (c *FakeBrucoProjects) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.BrucoProject, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(brucoprojectsResource, c.ns, name), &v1alpha1.BrucoProject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BrucoProject), err
}

// List takes label and field selectors, and returns the list of BrucoProjects that match those selectors.
func (c *FakeBrucoProjects) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BrucoProjectList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(brucoprojectsResource, brucoprojectsKind, c.ns, opts), &v1alpha1.BrucoProjectList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.BrucoProjectList{ListMeta: obj.(*v1alpha1.BrucoProjectList).ListMeta}
	for _, item := range obj.(*v1alpha1.BrucoProjectList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested brucoProjects.
func (c *FakeBrucoProjects) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(brucoprojectsResource, c.ns, opts))

}

// Create takes the representation of a brucoProject and creates it.  Returns the server's representation of the brucoProject, and an error, if there is any.
func (c *FakeBrucoProjects) Create(ctx context.Context, brucoProject *v1alpha1.BrucoProject, opts v1.CreateOptions) (result *v1alpha1.BrucoProject, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(brucoprojectsResource, c.ns, brucoProject), &v1alpha1.BrucoProject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BrucoProject), err
}

// Update takes the representation of a brucoProject and updates it. Returns the server's representation of the brucoProject, and an error, if there is any.
func (c *FakeBrucoProjects) Update(ctx context.Context, brucoProject *v1alpha1.BrucoProject, opts v1.UpdateOptions) (result *v1alpha1.BrucoProject, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(brucoprojectsResource, c.ns, brucoProject), &v1alpha1.BrucoProject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BrucoProject), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBrucoProjects) UpdateStatus(ctx context.Context, brucoProject *v1alpha1.BrucoProject, opts v1.UpdateOptions) (*v1alpha1.BrucoProject, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(brucoprojectsResource, "status", c.ns, brucoProject), &v1alpha1.BrucoProject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BrucoProject), err
}

// Delete takes name of the brucoProject and deletes it. Returns an error if one occurs.
func (c *FakeBrucoProjects) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(brucoprojectsResource, c.ns, name), &v1alpha1.BrucoProject{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBrucoProjects) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(brucoprojectsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.BrucoProjectList{})
	return err
}

// Patch applies the patch and returns the patched brucoProject.
func (c *FakeBrucoProjects) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.BrucoProject, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(brucoprojectsResource, c.ns, name, pt, data, subresources...), &v1alpha1.BrucoProject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.BrucoProject), err
}
