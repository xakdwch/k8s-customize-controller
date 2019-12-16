/*
Copyright 2019 The Kubernetes Authors.

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

package fake

import (
	bolingcavalry_v1 "k8s-customize-controller/pkg/apis/bolingcavalry/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeProgrammers implements ProgrammerInterface
type FakeProgrammers struct {
	Fake *FakeBolingcavalryV1
	ns   string
}

var programmersResource = schema.GroupVersionResource{Group: "bolingcavalry", Version: "v1", Resource: "programmers"}

var programmersKind = schema.GroupVersionKind{Group: "bolingcavalry", Version: "v1", Kind: "Programmer"}

// Get takes name of the programmer, and returns the corresponding programmer object, and an error if there is any.
func (c *FakeProgrammers) Get(name string, options v1.GetOptions) (result *bolingcavalry_v1.Programmer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(programmersResource, c.ns, name), &bolingcavalry_v1.Programmer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*bolingcavalry_v1.Programmer), err
}

// List takes label and field selectors, and returns the list of Programmers that match those selectors.
func (c *FakeProgrammers) List(opts v1.ListOptions) (result *bolingcavalry_v1.ProgrammerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(programmersResource, programmersKind, c.ns, opts), &bolingcavalry_v1.ProgrammerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &bolingcavalry_v1.ProgrammerList{}
	for _, item := range obj.(*bolingcavalry_v1.ProgrammerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested programmers.
func (c *FakeProgrammers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(programmersResource, c.ns, opts))

}

// Create takes the representation of a programmer and creates it.  Returns the server's representation of the programmer, and an error, if there is any.
func (c *FakeProgrammers) Create(programmer *bolingcavalry_v1.Programmer) (result *bolingcavalry_v1.Programmer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(programmersResource, c.ns, programmer), &bolingcavalry_v1.Programmer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*bolingcavalry_v1.Programmer), err
}

// Update takes the representation of a programmer and updates it. Returns the server's representation of the programmer, and an error, if there is any.
func (c *FakeProgrammers) Update(programmer *bolingcavalry_v1.Programmer) (result *bolingcavalry_v1.Programmer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(programmersResource, c.ns, programmer), &bolingcavalry_v1.Programmer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*bolingcavalry_v1.Programmer), err
}

// Delete takes name of the programmer and deletes it. Returns an error if one occurs.
func (c *FakeProgrammers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(programmersResource, c.ns, name), &bolingcavalry_v1.Programmer{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeProgrammers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(programmersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &bolingcavalry_v1.ProgrammerList{})
	return err
}

// Patch applies the patch and returns the patched programmer.
func (c *FakeProgrammers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *bolingcavalry_v1.Programmer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(programmersResource, c.ns, name, data, subresources...), &bolingcavalry_v1.Programmer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*bolingcavalry_v1.Programmer), err
}
