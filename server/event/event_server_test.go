package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	fakekube "k8s.io/client-go/kubernetes/fake"

	eventpkg "github.com/simster7/argo/v2/pkg/apiclient/event"
	wfv1 "github.com/simster7/argo/v2/pkg/apis/workflow/v1alpha1"
	"github.com/simster7/argo/v2/pkg/client/clientset/versioned/fake"
	"github.com/simster7/argo/v2/server/auth"
	"github.com/simster7/argo/v2/util/instanceid"
	"github.com/simster7/argo/v2/workflow/events"
)

func TestController(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	s := NewController(instanceid.NewService("my-instanceid"), events.NewEventRecorderManager(fakekube.NewSimpleClientset()), 1, 1)

	ctx := context.WithValue(context.TODO(), auth.WfKey, clientset)
	_, err := s.ReceiveEvent(ctx, &eventpkg.EventRequest{Namespace: "my-ns", Payload: &wfv1.Item{}})
	assert.NoError(t, err)

	assert.Len(t, s.operationQueue, 1, "one event to be processed")

	_, err = s.ReceiveEvent(ctx, &eventpkg.EventRequest{})
	assert.EqualError(t, err, "operation queue full", "backpressure when queue is full")

	stopCh := make(chan struct{}, 1)
	stopCh <- struct{}{}
	s.Run(stopCh)

	assert.Len(t, s.operationQueue, 0, "all events were processed")
}
