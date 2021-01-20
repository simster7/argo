package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/simster7/argo/v2/errors"
	wfv1 "github.com/simster7/argo/v2/pkg/apis/workflow/v1alpha1"
	"github.com/simster7/argo/v2/workflow/common"
)

// applyExecutionControl will ensure a pod's execution control annotation is up-to-date
// kills any pending pods when workflow has reached it's deadline
func (woc *wfOperationCtx) applyExecutionControl(ctx context.Context, pod *apiv1.Pod, wfNodesLock *sync.RWMutex) error {
	if pod == nil {
		return nil
	}
	switch pod.Status.Phase {
	case apiv1.PodSucceeded, apiv1.PodFailed:
		// Skip any pod which are already completed
		return nil
	case apiv1.PodPending:
		// Check if we are currently shutting down
		if woc.execWf.Spec.Shutdown != "" {
			// Only delete pods that are not part of an onExit handler if we are "Stopping" or all pods if we are "Terminating"
			_, onExitPod := pod.Labels[common.LabelKeyOnExit]

			if !woc.wf.Spec.Shutdown.ShouldExecute(onExitPod) {
				woc.log.Infof("Deleting Pending pod %s/%s as part of workflow shutdown with strategy: %s", pod.Namespace, pod.Name, woc.wf.Spec.Shutdown)
				err := woc.controller.kubeclientset.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
				if err == nil {
					wfNodesLock.Lock()
					defer wfNodesLock.Unlock()
					node := woc.wf.Status.Nodes[pod.Name]
					woc.markNodePhase(node.Name, wfv1.NodeFailed, fmt.Sprintf("workflow shutdown with strategy:  %s", woc.execWf.Spec.Shutdown))
					return nil
				}
				// If we fail to delete the pod, fall back to setting the annotation
				woc.log.Warnf("Failed to delete %s/%s: %v", pod.Namespace, pod.Name, err)
			}
		}
		// Check if we are past the workflow deadline. If we are, and the pod is still pending
		// then we should simply delete it and mark the pod as Failed
		if woc.workflowDeadline != nil && time.Now().UTC().After(*woc.workflowDeadline) {
			//pods that are part of an onExit handler aren't subject to the deadline
			_, onExitPod := pod.Labels[common.LabelKeyOnExit]
			if !onExitPod {
				woc.log.Infof("Deleting Pending pod %s/%s which has exceeded workflow deadline %s", pod.Namespace, pod.Name, woc.workflowDeadline)
				err := woc.controller.kubeclientset.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
				if err == nil {
					wfNodesLock.Lock()
					defer wfNodesLock.Unlock()
					node := woc.wf.Status.Nodes[pod.Name]
					woc.markNodePhase(node.Name, wfv1.NodeFailed, "Step exceeded its deadline")
					return nil
				}
				// If we fail to delete the pod, fall back to setting the annotation
				woc.log.Warnf("Failed to delete %s/%s: %v", pod.Namespace, pod.Name, err)
			}
		}
	}

	var podExecCtl common.ExecutionControl
	if execCtlStr, ok := pod.Annotations[common.AnnotationKeyExecutionControl]; ok && execCtlStr != "" {
		err := json.Unmarshal([]byte(execCtlStr), &podExecCtl)
		if err != nil {
			woc.log.Warnf("Failed to unmarshal execution control from pod %s", pod.Name)
		}
	}
	containerName := common.WaitContainerName
	// A resource template does not have a wait container,
	// instead the only container is the main container (which is running argoexec)
	if len(pod.Spec.Containers) == 1 {
		containerName = common.MainContainerName
	}

	if woc.wf.Spec.Shutdown != "" {
		if _, onExitPod := pod.Labels[common.LabelKeyOnExit]; !woc.wf.Spec.Shutdown.ShouldExecute(onExitPod) {
			podExecCtl.Deadline = &time.Time{}
			woc.log.Infof("Applying shutdown deadline for pod %s", pod.Name)
			return woc.updateExecutionControl(ctx, pod.Name, podExecCtl, containerName)
		}
	}

	if woc.workflowDeadline != nil {
		if podExecCtl.Deadline == nil || woc.workflowDeadline.Before(*podExecCtl.Deadline) {
			podExecCtl.Deadline = woc.workflowDeadline
			woc.log.Infof("Applying sooner Workflow Deadline for pod %s at: %v", pod.Name, woc.workflowDeadline)
			return woc.updateExecutionControl(ctx, pod.Name, podExecCtl, containerName)
		}
	}

	return nil
}

// killDaemonedChildren kill any daemoned pods of a steps or DAG template node.
func (woc *wfOperationCtx) killDaemonedChildren(ctx context.Context, nodeID string) error {
	woc.log.Infof("Checking daemoned children of %s", nodeID)
	var firstErr error
	execCtl := common.ExecutionControl{
		Deadline: &time.Time{},
	}
	for _, childNode := range woc.wf.Status.Nodes {
		if childNode.BoundaryID != nodeID {
			continue
		}
		if childNode.Daemoned == nil || !*childNode.Daemoned {
			continue
		}
		err := woc.updateExecutionControl(ctx, childNode.ID, execCtl, common.WaitContainerName)
		if err != nil {
			woc.log.Errorf("Failed to update execution control of node %s: %+v", childNode.ID, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// updateExecutionControl updates the execution control parameters
func (woc *wfOperationCtx) updateExecutionControl(ctx context.Context, podName string, execCtl common.ExecutionControl, containerName string) error {
	execCtlBytes, err := json.Marshal(execCtl)
	if err != nil {
		return errors.InternalWrapError(err)
	}

	woc.log.Infof("Updating execution control of %s: %s", podName, execCtlBytes)
	err = common.AddPodAnnotation(
		ctx,
		woc.controller.kubeclientset,
		podName,
		woc.wf.ObjectMeta.Namespace,
		common.AnnotationKeyExecutionControl,
		string(execCtlBytes),
	)
	if err != nil {
		return err
	}

	// Ideally we would simply annotate the pod with the updates and be done with it, allowing
	// the executor to notice the updates naturally via the Downward API annotations volume
	// mounted file. However, updates to the Downward API volumes take a very long time to
	// propagate (minutes). The following code fast-tracks this by signaling the executor
	// using SIGUSR2 that something changed.
	woc.log.Infof("Signalling %s of updates", podName)
	exec, err := common.ExecPodContainer(
		woc.controller.restConfig, woc.wf.ObjectMeta.Namespace, podName,
		containerName, true, true, "sh", "-c", "kill -s USR2 $(pidof argoexec)",
	)
	if err != nil {
		return err
	}
	go func() {
		// This call is necessary to actually send the exec. Since signalling is best effort,
		// it is launched as a goroutine and the error is discarded
		_, _, err = common.GetExecutorOutput(exec)
		if err != nil {
			woc.log.Warnf("Signal command failed: %v", err)
			return
		}
		woc.log.Infof("Signal of %s (%s) successfully issued", podName, common.WaitContainerName)
	}()

	return nil
}
