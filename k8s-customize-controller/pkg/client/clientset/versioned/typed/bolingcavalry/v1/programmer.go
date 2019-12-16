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

package v1

import (
	v1 "k8s-customize-controller/pkg/apis/bolingcavalry/v1"
	scheme "k8s-customize-controller/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ProgrammersGetter has a method to return a ProgrammerInterface.
// A group's client should implement this interface.
type ProgrammersGetter interface {
	Programmers(namespace string) ProgrammerInterface
}

// ProgrammerInterface has methods to work with Programmer resources.
type ProgrammerInterface interface {
	Create(*v1.Programmer) (*v1.Programmer, error)
	Update(*v1.Programmer) (*v1.Programmer, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Programmer, error)
	List(opts meta_v1.ListOptions) (*v1.ProgrammerList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Programmer, err error)
	ProgrammerExpansion
}

// programmers implements ProgrammerInterface
type programmers struct {
	client rest.Interface
	ns     string
}

// newProgrammers returns a Programmers
func newProgrammers(c *BolingcavalryV1Client, namespace string) *programmers {
	return &programmers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the programmer, and returns the corresponding programmer object, and an error if there is any.
func (c *programmers) Get(name string, options meta_v1.GetOptions) (result *v1.Programmer, err error) {
	result = &v1.Programmer{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("programmers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Programmers that match those selectors.
func (c *programmers) List(opts meta_v1.ListOptions) (result *v1.ProgrammerList, err error) {
	result = &v1.ProgrammerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("programmers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested programmers.
func (c *programmers) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("programmers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a programmer and creates it.  Returns the server's representation of the programmer, and an error, if there is any.
func (c *programmers) Create(programmer *v1.Programmer) (result *v1.Programmer, err error) {
	result = &v1.Programmer{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("programmers").
		Body(programmer).
		Do().
		Into(result)
	return
}

// Update takes the representation of a programmer and updates it. Returns the server's representation of the programmer, and an error, if there is any.
func (c *programmers) Update(programmer *v1.Programmer) (result *v1.Programmer, err error) {
	result = &v1.Programmer{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("programmers").
		Name(programmer.Name).
		Body(programmer).
		Do().
		Into(result)
	return
}

// Delete takes name of the programmer and deletes it. Returns an error if one occurs.
func (c *programmers) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("programmers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *programmers) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("programmers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched programmer.
func (c *programmers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Programmer, err error) {
	result = &v1.Programmer{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("programmers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
