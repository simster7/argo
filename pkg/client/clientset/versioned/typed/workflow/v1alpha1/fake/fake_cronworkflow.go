// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/simster7/argo/v2/pkg/apis/workflow/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCronWorkflows implements CronWorkflowInterface
type FakeCronWorkflows struct {
	Fake *FakeArgoprojV1alpha1
	ns   string
}

var cronworkflowsResource = schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "cronworkflows"}

var cronworkflowsKind = schema.GroupVersionKind{Group: "argoproj.io", Version: "v1alpha1", Kind: "CronWorkflow"}

// Get takes name of the cronWorkflow, and returns the corresponding cronWorkflow object, and an error if there is any.
func (c *FakeCronWorkflows) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.CronWorkflow, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(cronworkflowsResource, c.ns, name), &v1alpha1.CronWorkflow{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CronWorkflow), err
}

// List takes label and field selectors, and returns the list of CronWorkflows that match those selectors.
func (c *FakeCronWorkflows) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.CronWorkflowList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(cronworkflowsResource, cronworkflowsKind, c.ns, opts), &v1alpha1.CronWorkflowList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.CronWorkflowList{ListMeta: obj.(*v1alpha1.CronWorkflowList).ListMeta}
	for _, item := range obj.(*v1alpha1.CronWorkflowList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested cronWorkflows.
func (c *FakeCronWorkflows) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(cronworkflowsResource, c.ns, opts))

}

// Create takes the representation of a cronWorkflow and creates it.  Returns the server's representation of the cronWorkflow, and an error, if there is any.
func (c *FakeCronWorkflows) Create(ctx context.Context, cronWorkflow *v1alpha1.CronWorkflow, opts v1.CreateOptions) (result *v1alpha1.CronWorkflow, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(cronworkflowsResource, c.ns, cronWorkflow), &v1alpha1.CronWorkflow{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CronWorkflow), err
}

// Update takes the representation of a cronWorkflow and updates it. Returns the server's representation of the cronWorkflow, and an error, if there is any.
func (c *FakeCronWorkflows) Update(ctx context.Context, cronWorkflow *v1alpha1.CronWorkflow, opts v1.UpdateOptions) (result *v1alpha1.CronWorkflow, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(cronworkflowsResource, c.ns, cronWorkflow), &v1alpha1.CronWorkflow{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CronWorkflow), err
}

// Delete takes name of the cronWorkflow and deletes it. Returns an error if one occurs.
func (c *FakeCronWorkflows) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(cronworkflowsResource, c.ns, name), &v1alpha1.CronWorkflow{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCronWorkflows) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(cronworkflowsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.CronWorkflowList{})
	return err
}

// Patch applies the patch and returns the patched cronWorkflow.
func (c *FakeCronWorkflows) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.CronWorkflow, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(cronworkflowsResource, c.ns, name, pt, data, subresources...), &v1alpha1.CronWorkflow{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CronWorkflow), err
}
