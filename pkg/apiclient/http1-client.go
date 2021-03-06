package apiclient

import (
	"context"

	"github.com/simster7/argo/v2/pkg/apiclient/clusterworkflowtemplate"
	cronworkflowpkg "github.com/simster7/argo/v2/pkg/apiclient/cronworkflow"
	"github.com/simster7/argo/v2/pkg/apiclient/http1"
	infopkg "github.com/simster7/argo/v2/pkg/apiclient/info"
	workflowpkg "github.com/simster7/argo/v2/pkg/apiclient/workflow"
	workflowarchivepkg "github.com/simster7/argo/v2/pkg/apiclient/workflowarchive"
	workflowtemplatepkg "github.com/simster7/argo/v2/pkg/apiclient/workflowtemplate"
)

type httpClient http1.Facade

func (h httpClient) NewArchivedWorkflowServiceClient() (workflowarchivepkg.ArchivedWorkflowServiceClient, error) {
	return http1.ArchivedWorkflowsServiceClient(h), nil
}

func (h httpClient) NewWorkflowServiceClient() workflowpkg.WorkflowServiceClient {
	return http1.WorkflowServiceClient(h)
}

func (h httpClient) NewCronWorkflowServiceClient() cronworkflowpkg.CronWorkflowServiceClient {
	return http1.CronWorkflowServiceClient(h)
}

func (h httpClient) NewWorkflowTemplateServiceClient() workflowtemplatepkg.WorkflowTemplateServiceClient {
	return http1.WorkflowTemplateServiceClient(h)
}

func (h httpClient) NewClusterWorkflowTemplateServiceClient() clusterworkflowtemplate.ClusterWorkflowTemplateServiceClient {
	return http1.ClusterWorkflowTemplateServiceClient(h)
}

func (h httpClient) NewInfoServiceClient() (infopkg.InfoServiceClient, error) {
	return http1.InfoServiceClient(h), nil
}

func newHTTP1Client(baseUrl string, auth string) (context.Context, Client, error) {
	return context.Background(), httpClient(http1.NewFacade(baseUrl, auth)), nil
}
