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

// BrucoInformer provides access to a shared informer and lister for
// Brucos.
type BrucoInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.BrucoLister
}

type brucoInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewBrucoInformer constructs a new informer for Bruco type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBrucoInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBrucoInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredBrucoInformer constructs a new informer for Bruco type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBrucoInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BrucocontrollerV1alpha1().Brucos(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BrucocontrollerV1alpha1().Brucos(namespace).Watch(context.TODO(), options)
			},
		},
		&brucocontrollerv1alpha1.Bruco{},
		resyncPeriod,
		indexers,
	)
}

func (f *brucoInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBrucoInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *brucoInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&brucocontrollerv1alpha1.Bruco{}, f.defaultInformer)
}

func (f *brucoInformer) Lister() v1alpha1.BrucoLister {
	return v1alpha1.NewBrucoLister(f.Informer().GetIndexer())
}
