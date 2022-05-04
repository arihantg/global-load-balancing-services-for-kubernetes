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

package fake

import (
	v1alpha1 "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1alpha1/clientset/versioned/typed/ako/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAkoV1alpha1 struct {
	*testing.Fake
}

func (c *FakeAkoV1alpha1) AviInfraSettings() v1alpha1.AviInfraSettingInterface {
	return &FakeAviInfraSettings{c}
}

func (c *FakeAkoV1alpha1) ClusterSets(namespace string) v1alpha1.ClusterSetInterface {
	return &FakeClusterSets{c, namespace}
}

func (c *FakeAkoV1alpha1) HTTPRules(namespace string) v1alpha1.HTTPRuleInterface {
	return &FakeHTTPRules{c, namespace}
}

func (c *FakeAkoV1alpha1) HostRules(namespace string) v1alpha1.HostRuleInterface {
	return &FakeHostRules{c, namespace}
}

func (c *FakeAkoV1alpha1) MultiClusterIngresses(namespace string) v1alpha1.MultiClusterIngressInterface {
	return &FakeMultiClusterIngresses{c, namespace}
}

func (c *FakeAkoV1alpha1) ServiceImports(namespace string) v1alpha1.ServiceImportInterface {
	return &FakeServiceImports{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeAkoV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}