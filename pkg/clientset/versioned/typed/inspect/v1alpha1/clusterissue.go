// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/getupio-undistro/zora/apis/zora/v1alpha1"
	scheme "github.com/getupio-undistro/zora/pkg/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ClusterIssuesGetter has a method to return a ClusterIssueInterface.
// A group's client should implement this interface.
type ClusterIssuesGetter interface {
	ClusterIssues(namespace string) ClusterIssueInterface
}

// ClusterIssueInterface has methods to work with ClusterIssue resources.
type ClusterIssueInterface interface {
	Create(ctx context.Context, clusterIssue *v1alpha1.ClusterIssue, opts v1.CreateOptions) (*v1alpha1.ClusterIssue, error)
	Update(ctx context.Context, clusterIssue *v1alpha1.ClusterIssue, opts v1.UpdateOptions) (*v1alpha1.ClusterIssue, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ClusterIssue, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ClusterIssueList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterIssue, err error)
	ClusterIssueExpansion
}

// clusterIssues implements ClusterIssueInterface
type clusterIssues struct {
	client rest.Interface
	ns     string
}

// newClusterIssues returns a ClusterIssues
func newClusterIssues(c *ZoraV1alpha1Client, namespace string) *clusterIssues {
	return &clusterIssues{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the clusterIssue, and returns the corresponding clusterIssue object, and an error if there is any.
func (c *clusterIssues) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ClusterIssue, err error) {
	result = &v1alpha1.ClusterIssue{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("clusterissues").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterIssues that match those selectors.
func (c *clusterIssues) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ClusterIssueList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ClusterIssueList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("clusterissues").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterIssues.
func (c *clusterIssues) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("clusterissues").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a clusterIssue and creates it.  Returns the server's representation of the clusterIssue, and an error, if there is any.
func (c *clusterIssues) Create(ctx context.Context, clusterIssue *v1alpha1.ClusterIssue, opts v1.CreateOptions) (result *v1alpha1.ClusterIssue, err error) {
	result = &v1alpha1.ClusterIssue{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("clusterissues").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterIssue).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a clusterIssue and updates it. Returns the server's representation of the clusterIssue, and an error, if there is any.
func (c *clusterIssues) Update(ctx context.Context, clusterIssue *v1alpha1.ClusterIssue, opts v1.UpdateOptions) (result *v1alpha1.ClusterIssue, err error) {
	result = &v1alpha1.ClusterIssue{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("clusterissues").
		Name(clusterIssue.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterIssue).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the clusterIssue and deletes it. Returns an error if one occurs.
func (c *clusterIssues) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("clusterissues").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterIssues) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("clusterissues").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched clusterIssue.
func (c *clusterIssues) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterIssue, err error) {
	result = &v1alpha1.ClusterIssue{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("clusterissues").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
