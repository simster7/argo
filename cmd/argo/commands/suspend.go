package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/simster7/argo/v2/cmd/argo/commands/client"
	workflowpkg "github.com/simster7/argo/v2/pkg/apiclient/workflow"
)

func NewSuspendCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "suspend WORKFLOW1 WORKFLOW2...",
		Short: "suspend zero or more workflow",
		Example: `# Suspend a workflow:

  argo suspend my-wf

# Suspend the latest workflow:
  argo suspend @latest
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, apiClient := client.NewAPIClient()
			serviceClient := apiClient.NewWorkflowServiceClient()
			namespace := client.Namespace()
			for _, wfName := range args {
				_, err := serviceClient.SuspendWorkflow(ctx, &workflowpkg.WorkflowSuspendRequest{
					Name:      wfName,
					Namespace: namespace,
				})
				if err != nil {
					log.Fatalf("Failed to suspended %s: %+v", wfName, err)
				}
				fmt.Printf("workflow %s suspended\n", wfName)
			}
		},
	}
	return command
}
