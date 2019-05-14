/*
Copyright 2019 Rancher Labs.

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

package v1alpha3

import (
	"context"

	v1alpha3 "github.com/knative/pkg/apis/istio/v1alpha3"
	clientset "github.com/knative/pkg/client/clientset/versioned/typed/istio/v1alpha3"
	informers "github.com/knative/pkg/client/informers/externalversions/istio/v1alpha3"
	listers "github.com/knative/pkg/client/listers/istio/v1alpha3"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type GatewayHandler func(string, *v1alpha3.Gateway) (*v1alpha3.Gateway, error)

type GatewayController interface {
	GatewayClient

	OnChange(ctx context.Context, name string, sync GatewayHandler)
	OnRemove(ctx context.Context, name string, sync GatewayHandler)
	Enqueue(namespace, name string)

	Cache() GatewayCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type GatewayClient interface {
	Create(*v1alpha3.Gateway) (*v1alpha3.Gateway, error)
	Update(*v1alpha3.Gateway) (*v1alpha3.Gateway, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha3.Gateway, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha3.GatewayList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha3.Gateway, err error)
}

type GatewayCache interface {
	Get(namespace, name string) (*v1alpha3.Gateway, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha3.Gateway, error)

	AddIndexer(indexName string, indexer GatewayIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha3.Gateway, error)
}

type GatewayIndexer func(obj *v1alpha3.Gateway) ([]string, error)

type gatewayController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.GatewaysGetter
	informer          informers.GatewayInformer
	gvk               schema.GroupVersionKind
}

func NewGatewayController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.GatewaysGetter, informer informers.GatewayInformer) GatewayController {
	return &gatewayController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromGatewayHandlerToHandler(sync GatewayHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha3.Gateway
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha3.Gateway))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *gatewayController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha3.Gateway))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateGatewayOnChange(updater generic.Updater, handler GatewayHandler) GatewayHandler {
	return func(key string, obj *v1alpha3.Gateway) (*v1alpha3.Gateway, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1alpha3.Gateway)
			}
		}

		return copyObj, err
	}
}

func (c *gatewayController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *gatewayController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *gatewayController) OnChange(ctx context.Context, name string, sync GatewayHandler) {
	c.AddGenericHandler(ctx, name, FromGatewayHandlerToHandler(sync))
}

func (c *gatewayController) OnRemove(ctx context.Context, name string, sync GatewayHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromGatewayHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *gatewayController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, namespace, name)
}

func (c *gatewayController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *gatewayController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *gatewayController) Cache() GatewayCache {
	return &gatewayCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *gatewayController) Create(obj *v1alpha3.Gateway) (*v1alpha3.Gateway, error) {
	return c.clientGetter.Gateways(obj.Namespace).Create(obj)
}

func (c *gatewayController) Update(obj *v1alpha3.Gateway) (*v1alpha3.Gateway, error) {
	return c.clientGetter.Gateways(obj.Namespace).Update(obj)
}

func (c *gatewayController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Gateways(namespace).Delete(name, options)
}

func (c *gatewayController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha3.Gateway, error) {
	return c.clientGetter.Gateways(namespace).Get(name, options)
}

func (c *gatewayController) List(namespace string, opts metav1.ListOptions) (*v1alpha3.GatewayList, error) {
	return c.clientGetter.Gateways(namespace).List(opts)
}

func (c *gatewayController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Gateways(namespace).Watch(opts)
}

func (c *gatewayController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha3.Gateway, err error) {
	return c.clientGetter.Gateways(namespace).Patch(name, pt, data, subresources...)
}

type gatewayCache struct {
	lister  listers.GatewayLister
	indexer cache.Indexer
}

func (c *gatewayCache) Get(namespace, name string) (*v1alpha3.Gateway, error) {
	return c.lister.Gateways(namespace).Get(name)
}

func (c *gatewayCache) List(namespace string, selector labels.Selector) ([]*v1alpha3.Gateway, error) {
	return c.lister.Gateways(namespace).List(selector)
}

func (c *gatewayCache) AddIndexer(indexName string, indexer GatewayIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha3.Gateway))
		},
	}))
}

func (c *gatewayCache) GetByIndex(indexName, key string) (result []*v1alpha3.Gateway, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha3.Gateway))
	}
	return result, nil
}