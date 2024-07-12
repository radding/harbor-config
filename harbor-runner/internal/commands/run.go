package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(RunCommand)
}

var RunCommand = &cobra.Command{
	Use:   "run",
	Short: "Run a task in the harbor workspace/project",
	Long: `Run a registered task in the harbor project or Workspace.
	If this command is run in a workspace, Harbor will go through all projects and find tasks with the same name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("I will run here!")
		return nil
	},
}
