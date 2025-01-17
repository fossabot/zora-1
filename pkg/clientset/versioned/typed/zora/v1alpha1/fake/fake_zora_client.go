// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/getupio-undistro/zora/pkg/clientset/versioned/typed/zora/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeZoraV1alpha1 struct {
	*testing.Fake
}

func (c *FakeZoraV1alpha1) Clusters(namespace string) v1alpha1.ClusterInterface {
	return &FakeClusters{c, namespace}
}

func (c *FakeZoraV1alpha1) ClusterIssues(namespace string) v1alpha1.ClusterIssueInterface {
	return &FakeClusterIssues{c, namespace}
}

func (c *FakeZoraV1alpha1) ClusterScans(namespace string) v1alpha1.ClusterScanInterface {
	return &FakeClusterScans{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeZoraV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
