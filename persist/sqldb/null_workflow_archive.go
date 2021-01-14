package sqldb

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"

	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

var NullWorkflowArchive WorkflowArchive = &nullWorkflowArchive{}

type nullWorkflowArchive struct {
}

func (r *nullWorkflowArchive) IsEnabled() bool {
	return false
}

func (r *nullWorkflowArchive) ArchiveWorkflow(wf *wfv1.Workflow) error {
	log.WithFields(log.Fields{"namespace": wf.Namespace, "workflow": wf, "uid": wf.UID}).Info("Logging workflow in lieu of archiving")
	return nil
}

func (r *nullWorkflowArchive) ListWorkflows(string, time.Time, time.Time, labels.Requirements, int, int) (wfv1.Workflows, error) {
	return wfv1.Workflows{}, nil
}

func (r *nullWorkflowArchive) GetWorkflow(string) (*wfv1.Workflow, error) {
	return nil, fmt.Errorf("getting archived workflows not supported")
}

func (r *nullWorkflowArchive) DeleteWorkflow(string) error {
	return fmt.Errorf("deleting archived workflows not supported")
}

func (r *nullWorkflowArchive) DeleteExpiredWorkflows(time.Duration) error {
	return nil
}
