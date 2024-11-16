package commands

import (
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	packageconfig "github.com/radding/harbor-runner/internal/package-config"
	"github.com/radding/harbor-runner/internal/taskgraph"
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
		cfg := packageconfig.GetConfig()
		if cfg == nil {
			return errors.New("failed to run command, no configuration found")
		}
		ctx := cfg.ConfigureContext(cmd.Context())
		tree, err := taskgraph.CreateTreeFromConfig(cfg, exec)
		if err != nil {
			return errors.Wrap(err, "failed to build task tree")
		}
		if !cfg.WasSetupRun {
			slog.Debug("Looks like this setup was never run, running setup now")
			fmt.Println("setting up package")
			err := tree.RunSetup(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to run package setup")
			}
			cfg.WasSetupRun = true
			cfg.Save()
		}
		tree.RunTask(ctx, args[0])
		return nil
	},
}
