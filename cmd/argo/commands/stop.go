package commands

import (
	"fmt"

	"github.com/argoproj/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/argoproj/argo/cmd/argo/commands/client"
	workflowpkg "github.com/argoproj/argo/pkg/apiclient/workflow"
)

func NewStopCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "stop WORKFLOW WORKFLOW2...",
		Short: "stop zero or more workflows",
		Run: func(cmd *cobra.Command, args []string) {

			ctx, apiClient := client.NewAPIClient()
			serviceClient := apiClient.NewWorkflowServiceClient()
			namespace := client.Namespace()
			for _, name := range args {
				wf, err := serviceClient.StopWorkflow(ctx, &workflowpkg.WorkflowStopRequest{
					Name:      name,
					Namespace: namespace,
				})
				errors.CheckError(err)
				fmt.Printf("workflow %s stopped\n", wf.Name)
			}
		},
	}
	return command
}
