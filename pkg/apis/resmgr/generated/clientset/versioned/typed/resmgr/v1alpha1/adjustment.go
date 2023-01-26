// Copyright 2019-2020 Intel Corporation. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	scheme "github.com/intel/nri-resmgr/pkg/apis/resmgr/generated/clientset/versioned/scheme"
	v1alpha1 "github.com/intel/nri-resmgr/pkg/apis/resmgr/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AdjustmentsGetter has a method to return a AdjustmentInterface.
// A group's client should implement this interface.
type AdjustmentsGetter interface {
	Adjustments(namespace string) AdjustmentInterface
}

// AdjustmentInterface has methods to work with Adjustment resources.
type AdjustmentInterface interface {
	Create(ctx context.Context, adjustment *v1alpha1.Adjustment, opts v1.CreateOptions) (*v1alpha1.Adjustment, error)
	Update(ctx context.Context, adjustment *v1alpha1.Adjustment, opts v1.UpdateOptions) (*v1alpha1.Adjustment, error)
	UpdateStatus(ctx context.Context, adjustment *v1alpha1.Adjustment, opts v1.UpdateOptions) (*v1alpha1.Adjustment, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Adjustment, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.AdjustmentList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Adjustment, err error)
	AdjustmentExpansion
}

// adjustments implements AdjustmentInterface
type adjustments struct {
	client rest.Interface
	ns     string
}

// newAdjustments returns a Adjustments
func newAdjustments(c *NriresmgrV1alpha1Client, namespace string) *adjustments {
	return &adjustments{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the adjustment, and returns the corresponding adjustment object, and an error if there is any.
func (c *adjustments) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Adjustment, err error) {
	result = &v1alpha1.Adjustment{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("adjustments").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Adjustments that match those selectors.
func (c *adjustments) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AdjustmentList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.AdjustmentList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("adjustments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested adjustments.
func (c *adjustments) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("adjustments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a adjustment and creates it.  Returns the server's representation of the adjustment, and an error, if there is any.
func (c *adjustments) Create(ctx context.Context, adjustment *v1alpha1.Adjustment, opts v1.CreateOptions) (result *v1alpha1.Adjustment, err error) {
	result = &v1alpha1.Adjustment{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("adjustments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(adjustment).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a adjustment and updates it. Returns the server's representation of the adjustment, and an error, if there is any.
func (c *adjustments) Update(ctx context.Context, adjustment *v1alpha1.Adjustment, opts v1.UpdateOptions) (result *v1alpha1.Adjustment, err error) {
	result = &v1alpha1.Adjustment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("adjustments").
		Name(adjustment.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(adjustment).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *adjustments) UpdateStatus(ctx context.Context, adjustment *v1alpha1.Adjustment, opts v1.UpdateOptions) (result *v1alpha1.Adjustment, err error) {
	result = &v1alpha1.Adjustment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("adjustments").
		Name(adjustment.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(adjustment).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the adjustment and deletes it. Returns an error if one occurs.
func (c *adjustments) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("adjustments").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *adjustments) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("adjustments").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched adjustment.
func (c *adjustments) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Adjustment, err error) {
	result = &v1alpha1.Adjustment{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("adjustments").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
