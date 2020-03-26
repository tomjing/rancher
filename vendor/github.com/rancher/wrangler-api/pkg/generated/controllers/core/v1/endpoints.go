/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/generic"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes/typed/core/v1"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type EndpointsHandler func(string, *v1.Endpoints) (*v1.Endpoints, error)

type EndpointsController interface {
	generic.ControllerMeta
	EndpointsClient

	OnChange(ctx context.Context, name string, sync EndpointsHandler)
	OnRemove(ctx context.Context, name string, sync EndpointsHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() EndpointsCache
}

type EndpointsClient interface {
	Create(*v1.Endpoints) (*v1.Endpoints, error)
	Update(*v1.Endpoints) (*v1.Endpoints, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Endpoints, error)
	List(namespace string, opts metav1.ListOptions) (*v1.EndpointsList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Endpoints, err error)
}

type EndpointsCache interface {
	Get(namespace, name string) (*v1.Endpoints, error)
	List(namespace string, selector labels.Selector) ([]*v1.Endpoints, error)

	AddIndexer(indexName string, indexer EndpointsIndexer)
	GetByIndex(indexName, key string) ([]*v1.Endpoints, error)
}

type EndpointsIndexer func(obj *v1.Endpoints) ([]string, error)

type endpointsController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.EndpointsGetter
	informer          informers.EndpointsInformer
	gvk               schema.GroupVersionKind
}

func NewEndpointsController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.EndpointsGetter, informer informers.EndpointsInformer) EndpointsController {
	return &endpointsController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromEndpointsHandlerToHandler(sync EndpointsHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Endpoints
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Endpoints))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *endpointsController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Endpoints))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateEndpointsDeepCopyOnChange(client EndpointsClient, obj *v1.Endpoints, handler func(obj *v1.Endpoints) (*v1.Endpoints, error)) (*v1.Endpoints, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *endpointsController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *endpointsController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *endpointsController) OnChange(ctx context.Context, name string, sync EndpointsHandler) {
	c.AddGenericHandler(ctx, name, FromEndpointsHandlerToHandler(sync))
}

func (c *endpointsController) OnRemove(ctx context.Context, name string, sync EndpointsHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromEndpointsHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *endpointsController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *endpointsController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), namespace, name, duration)
}

func (c *endpointsController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *endpointsController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *endpointsController) Cache() EndpointsCache {
	return &endpointsCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *endpointsController) Create(obj *v1.Endpoints) (*v1.Endpoints, error) {
	return c.clientGetter.Endpoints(obj.Namespace).Create(obj)
}

func (c *endpointsController) Update(obj *v1.Endpoints) (*v1.Endpoints, error) {
	return c.clientGetter.Endpoints(obj.Namespace).Update(obj)
}

func (c *endpointsController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Endpoints(namespace).Delete(name, options)
}

func (c *endpointsController) Get(namespace, name string, options metav1.GetOptions) (*v1.Endpoints, error) {
	return c.clientGetter.Endpoints(namespace).Get(name, options)
}

func (c *endpointsController) List(namespace string, opts metav1.ListOptions) (*v1.EndpointsList, error) {
	return c.clientGetter.Endpoints(namespace).List(opts)
}

func (c *endpointsController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Endpoints(namespace).Watch(opts)
}

func (c *endpointsController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Endpoints, err error) {
	return c.clientGetter.Endpoints(namespace).Patch(name, pt, data, subresources...)
}

type endpointsCache struct {
	lister  listers.EndpointsLister
	indexer cache.Indexer
}

func (c *endpointsCache) Get(namespace, name string) (*v1.Endpoints, error) {
	return c.lister.Endpoints(namespace).Get(name)
}

func (c *endpointsCache) List(namespace string, selector labels.Selector) ([]*v1.Endpoints, error) {
	return c.lister.Endpoints(namespace).List(selector)
}

func (c *endpointsCache) AddIndexer(indexName string, indexer EndpointsIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Endpoints))
		},
	}))
}

func (c *endpointsCache) GetByIndex(indexName, key string) (result []*v1.Endpoints, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.Endpoints, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.Endpoints))
	}
	return result, nil
}