package http1

import (
	workflowpkg "github.com/simster7/argo/v2/pkg/apiclient/workflow"
)

type podLogsClient struct{ serverSentEventsClient }

func (f *podLogsClient) Recv() (*workflowpkg.LogEntry, error) {
	v := &workflowpkg.LogEntry{}
	return v, f.RecvEvent(v)
}
