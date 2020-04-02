package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourceScope string

const (
	ResourceScopeLocal                   ResourceScope = "Local"
	ResourceScopeWorkflowTemplate        ResourceScope = "WorkflowTemplate"
	ResourceScopeClusterWorkflowTemplate ResourceScope = "ClusterWorkflowTemplate"
)

// TemplateHolder is an object that holds templates, such as a Workflow, WorkflowTempalte, or ClusterWorkflowTemplate
type TemplateHolder interface {
	GetNamespace() string
	GetName() string
	GroupVersionKind() schema.GroupVersionKind
	GetTemplateByName(name string) *Template
	GetTemplateScope() (ResourceScope, string)
	GetAllTemplates() []Template
}

// TemplateCaller is an object that can call other templates, such as a Workflow, DAGTask, or WorkflowStep
type TemplateCaller interface {
	GetTemplateName() string
	GetTemplateRef() *TemplateRef
}
