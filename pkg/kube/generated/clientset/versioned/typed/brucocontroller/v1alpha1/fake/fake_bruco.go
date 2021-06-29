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

// FakeBrucos implements BrucoInterface
type FakeBrucos struct {
	Fake *FakeBrucoV1alpha1
	ns   string
}

var brucosResource = schema.GroupVersionResource{Group: "bruco.ferama.github.io", Version: "v1alpha1", Resource: "brucos"}

var brucosKind = schema.GroupVersionKind{Group: "bruco.ferama.github.io", Version: "v1alpha1", Kind: "Bruco"}

// Get takes name of the bruco, and returns the corresponding bruco object, and an error if there is any.
func (c *FakeBrucos) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Bruco, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(brucosResource, c.ns, name), &v1alpha1.Bruco{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bruco), err
}

// List takes label and field selectors, and returns the list of Brucos that match those selectors.
func (c *FakeBrucos) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BrucoList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(brucosResource, brucosKind, c.ns, opts), &v1alpha1.BrucoList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.BrucoList{ListMeta: obj.(*v1alpha1.BrucoList).ListMeta}
	for _, item := range obj.(*v1alpha1.BrucoList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested brucos.
func (c *FakeBrucos) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(brucosResource, c.ns, opts))

}

// Create takes the representation of a bruco and creates it.  Returns the server's representation of the bruco, and an error, if there is any.
func (c *FakeBrucos) Create(ctx context.Context, bruco *v1alpha1.Bruco, opts v1.CreateOptions) (result *v1alpha1.Bruco, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(brucosResource, c.ns, bruco), &v1alpha1.Bruco{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bruco), err
}

// Update takes the representation of a bruco and updates it. Returns the server's representation of the bruco, and an error, if there is any.
func (c *FakeBrucos) Update(ctx context.Context, bruco *v1alpha1.Bruco, opts v1.UpdateOptions) (result *v1alpha1.Bruco, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(brucosResource, c.ns, bruco), &v1alpha1.Bruco{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bruco), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBrucos) UpdateStatus(ctx context.Context, bruco *v1alpha1.Bruco, opts v1.UpdateOptions) (*v1alpha1.Bruco, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(brucosResource, "status", c.ns, bruco), &v1alpha1.Bruco{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bruco), err
}

// Delete takes name of the bruco and deletes it. Returns an error if one occurs.
func (c *FakeBrucos) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(brucosResource, c.ns, name), &v1alpha1.Bruco{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBrucos) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(brucosResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.BrucoList{})
	return err
}

// Patch applies the patch and returns the patched bruco.
func (c *FakeBrucos) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Bruco, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(brucosResource, c.ns, name, pt, data, subresources...), &v1alpha1.Bruco{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bruco), err
}