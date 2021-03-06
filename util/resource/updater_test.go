package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wfv1 "github.com/simster7/argo/v2/pkg/apis/workflow/v1alpha1"
	"github.com/simster7/argo/v2/test/util"
)

func TestUpdater(t *testing.T) {
	wf := &wfv1.Workflow{}
	util.MustUnmarshallYAML(`
status:
  nodes:
    root:
      phase: Succeeded
      children: [pod, dag] 
    pod: 
      phase: Succeeded
      type: Pod
      resourcesDuration: 
        x: 1
      children: [dag]
    dag: 
      phase: Succeeded
      children: [dag-pod]
    dag-pod: 
      phase: Succeeded
      type: Pod
      resourcesDuration: 
        x: 2
`, wf)
	UpdateResourceDurations(wf)
	assert.Equal(t, wfv1.ResourcesDuration{"x": 2}, wf.Status.Nodes["dag-pod"].ResourcesDuration)
	assert.Equal(t, wfv1.ResourcesDuration{"x": 2}, wf.Status.Nodes["dag"].ResourcesDuration)
	assert.Equal(t, wfv1.ResourcesDuration{"x": 1}, wf.Status.Nodes["pod"].ResourcesDuration)
	assert.Equal(t, wfv1.ResourcesDuration{"x": 3}, wf.Status.Nodes["root"].ResourcesDuration)
	assert.Equal(t, wfv1.ResourcesDuration{"x": 3}, wf.Status.ResourcesDuration)
}
