// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	brucocontrollerv1alpha1 "github.com/ferama/bruco/pkg/kube/apis/brucocontroller/v1alpha1"
	versioned "github.com/ferama/bruco/pkg/kube/generated/clientset/versioned"
	internalinterfaces "github.com/ferama/bruco/pkg/kube/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/ferama/bruco/pkg/kube/generated/listers/brucocontroller/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// BrucoProjectInformer provides access to a shared informer and lister for
// BrucoProjects.
type BrucoProjectInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.BrucoProjectLister
}

type brucoProjectInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewBrucoProjectInformer constructs a new informer for BrucoProject type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBrucoProjectInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBrucoProjectInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredBrucoProjectInformer constructs a new informer for BrucoProject type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBrucoProjectInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BrucoV1alpha1().BrucoProjects(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BrucoV1alpha1().BrucoProjects(namespace).Watch(context.TODO(), options)
			},
		},
		&brucocontrollerv1alpha1.BrucoProject{},
		resyncPeriod,
		indexers,
	)
}

func (f *brucoProjectInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBrucoProjectInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *brucoProjectInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&brucocontrollerv1alpha1.BrucoProject{}, f.defaultInformer)
}

func (f *brucoProjectInformer) Lister() v1alpha1.BrucoProjectLister {
	return v1alpha1.NewBrucoProjectLister(f.Informer().GetIndexer())
}
