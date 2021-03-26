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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/apis/ako/v1alpha1"
	scheme "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1alpha1/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HostRulesGetter has a method to return a HostRuleInterface.
// A group's client should implement this interface.
type HostRulesGetter interface {
	HostRules(namespace string) HostRuleInterface
}

// HostRuleInterface has methods to work with HostRule resources.
type HostRuleInterface interface {
	Create(ctx context.Context, hostRule *v1alpha1.HostRule, opts v1.CreateOptions) (*v1alpha1.HostRule, error)
	Update(ctx context.Context, hostRule *v1alpha1.HostRule, opts v1.UpdateOptions) (*v1alpha1.HostRule, error)
	UpdateStatus(ctx context.Context, hostRule *v1alpha1.HostRule, opts v1.UpdateOptions) (*v1alpha1.HostRule, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.HostRule, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.HostRuleList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.HostRule, err error)
	HostRuleExpansion
}

// hostRules implements HostRuleInterface
type hostRules struct {
	client rest.Interface
	ns     string
}

// newHostRules returns a HostRules
func newHostRules(c *AkoV1alpha1Client, namespace string) *hostRules {
	return &hostRules{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the hostRule, and returns the corresponding hostRule object, and an error if there is any.
func (c *hostRules) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.HostRule, err error) {
	result = &v1alpha1.HostRule{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hostrules").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of HostRules that match those selectors.
func (c *hostRules) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.HostRuleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.HostRuleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hostrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested hostRules.
func (c *hostRules) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("hostrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a hostRule and creates it.  Returns the server's representation of the hostRule, and an error, if there is any.
func (c *hostRules) Create(ctx context.Context, hostRule *v1alpha1.HostRule, opts v1.CreateOptions) (result *v1alpha1.HostRule, err error) {
	result = &v1alpha1.HostRule{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("hostrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hostRule).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a hostRule and updates it. Returns the server's representation of the hostRule, and an error, if there is any.
func (c *hostRules) Update(ctx context.Context, hostRule *v1alpha1.HostRule, opts v1.UpdateOptions) (result *v1alpha1.HostRule, err error) {
	result = &v1alpha1.HostRule{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hostrules").
		Name(hostRule.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hostRule).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *hostRules) UpdateStatus(ctx context.Context, hostRule *v1alpha1.HostRule, opts v1.UpdateOptions) (result *v1alpha1.HostRule, err error) {
	result = &v1alpha1.HostRule{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hostrules").
		Name(hostRule.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hostRule).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the hostRule and deletes it. Returns an error if one occurs.
func (c *hostRules) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hostrules").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *hostRules) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hostrules").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched hostRule.
func (c *hostRules) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.HostRule, err error) {
	result = &v1alpha1.HostRule{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("hostrules").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
