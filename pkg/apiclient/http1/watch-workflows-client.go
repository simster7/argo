package http1

import (
	workflowpkg "github.com/simster7/argo/v2/pkg/apiclient/workflow"
)

type watchWorkflowsClient struct{ serverSentEventsClient }

func (f watchWorkflowsClient) Recv() (*workflowpkg.WorkflowWatchEvent, error) {
	v := &workflowpkg.WorkflowWatchEvent{}
	return v, f.RecvEvent(v)
}
