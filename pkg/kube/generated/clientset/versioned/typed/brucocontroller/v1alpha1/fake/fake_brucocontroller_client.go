// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/ferama/bruco/pkg/kube/generated/clientset/versioned/typed/brucocontroller/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeBrucoV1alpha1 struct {
	*testing.Fake
}

func (c *FakeBrucoV1alpha1) Brucos(namespace string) v1alpha1.BrucoInterface {
	return &FakeBrucos{c, namespace}
}

func (c *FakeBrucoV1alpha1) BrucoProjects(namespace string) v1alpha1.BrucoProjectInterface {
	return &FakeBrucoProjects{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeBrucoV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
