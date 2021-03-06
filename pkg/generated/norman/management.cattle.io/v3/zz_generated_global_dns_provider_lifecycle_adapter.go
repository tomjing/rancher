package v3

import (
	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type GlobalDNSProviderLifecycle interface {
	Create(obj *v3.GlobalDNSProvider) (runtime.Object, error)
	Remove(obj *v3.GlobalDNSProvider) (runtime.Object, error)
	Updated(obj *v3.GlobalDNSProvider) (runtime.Object, error)
}

type globalDnsProviderLifecycleAdapter struct {
	lifecycle GlobalDNSProviderLifecycle
}

func (w *globalDnsProviderLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *globalDnsProviderLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *globalDnsProviderLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v3.GlobalDNSProvider))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalDnsProviderLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v3.GlobalDNSProvider))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalDnsProviderLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v3.GlobalDNSProvider))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGlobalDNSProviderLifecycleAdapter(name string, clusterScoped bool, client GlobalDNSProviderInterface, l GlobalDNSProviderLifecycle) GlobalDNSProviderHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(GlobalDNSProviderGroupVersionResource)
	}
	adapter := &globalDnsProviderLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.GlobalDNSProvider) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
