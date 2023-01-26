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

package fake

import (
	"context"

	v1alpha1 "github.com/intel/nri-resmgr/pkg/apis/resmgr/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAdjustments implements AdjustmentInterface
type FakeAdjustments struct {
	Fake *FakeNriresmgrV1alpha1
	ns   string
}

var adjustmentsResource = schema.GroupVersionResource{Group: "nriresmgr.intel.com", Version: "v1alpha1", Resource: "adjustments"}

var adjustmentsKind = schema.GroupVersionKind{Group: "nriresmgr.intel.com", Version: "v1alpha1", Kind: "Adjustment"}

// Get takes name of the adjustment, and returns the corresponding adjustment object, and an error if there is any.
func (c *FakeAdjustments) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Adjustment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(adjustmentsResource, c.ns, name), &v1alpha1.Adjustment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Adjustment), err
}

// List takes label and field selectors, and returns the list of Adjustments that match those selectors.
func (c *FakeAdjustments) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AdjustmentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(adjustmentsResource, adjustmentsKind, c.ns, opts), &v1alpha1.AdjustmentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AdjustmentList{ListMeta: obj.(*v1alpha1.AdjustmentList).ListMeta}
	for _, item := range obj.(*v1alpha1.AdjustmentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested adjustments.
func (c *FakeAdjustments) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(adjustmentsResource, c.ns, opts))

}

// Create takes the representation of a adjustment and creates it.  Returns the server's representation of the adjustment, and an error, if there is any.
func (c *FakeAdjustments) Create(ctx context.Context, adjustment *v1alpha1.Adjustment, opts v1.CreateOptions) (result *v1alpha1.Adjustment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(adjustmentsResource, c.ns, adjustment), &v1alpha1.Adjustment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Adjustment), err
}

// Update takes the representation of a adjustment and updates it. Returns the server's representation of the adjustment, and an error, if there is any.
func (c *FakeAdjustments) Update(ctx context.Context, adjustment *v1alpha1.Adjustment, opts v1.UpdateOptions) (result *v1alpha1.Adjustment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(adjustmentsResource, c.ns, adjustment), &v1alpha1.Adjustment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Adjustment), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAdjustments) UpdateStatus(ctx context.Context, adjustment *v1alpha1.Adjustment, opts v1.UpdateOptions) (*v1alpha1.Adjustment, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(adjustmentsResource, "status", c.ns, adjustment), &v1alpha1.Adjustment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Adjustment), err
}

// Delete takes name of the adjustment and deletes it. Returns an error if one occurs.
func (c *FakeAdjustments) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(adjustmentsResource, c.ns, name, opts), &v1alpha1.Adjustment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAdjustments) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(adjustmentsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.AdjustmentList{})
	return err
}

// Patch applies the patch and returns the patched adjustment.
func (c *FakeAdjustments) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Adjustment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(adjustmentsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Adjustment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Adjustment), err
}
