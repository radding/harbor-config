package commands

import (
	"github.com/spf13/cobra"
)

var force = false

func init() {
	rootCmd.AddCommand(SetupCommand)
	SetupCommand.Flags().BoolVarP(&force, "force", "f", false, "Force setup to run, even if not needed")
}

var SetupCommand = &cobra.Command{
	Use:   "setup",
	Short: "Run the package's or workspace's Package setup.",
	Long:  "Run the package's or workspace's Package setup. Setup is implicitly run when a task is run if needed",
	// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
	// _, err := harborconfig.LoadConfig("./.harborrc.ts")
	// return err
	// },
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
