package commands

import (
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	packageconfig "github.com/radding/harbor-runner/internal/package-config"
	"github.com/radding/harbor-runner/internal/taskgraph"
	"github.com/spf13/cobra"
)

func createSetupCommand(root *cobra.Command, exec taskgraph.Executor) {
	force := false

	SetupCommand := &cobra.Command{
		Use:   "setup",
		Short: "Run the package's or workspace's Package setup.",
		Long:  "Run the package's or workspace's Package setup. Setup is implicitly run when a task is run if needed",
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// _, err := harborconfig.LoadConfig("./.harborrc.ts")
		// return err
		// },
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := packageconfig.GetConfig()
			ctx := cfg.ConfigureContext(cmd.Context())
			tree, err := taskgraph.CreateTreeFromConfig(cfg, exec)
			if err != nil {
				return errors.Wrap(err, "failed to build task tree")
			}
			if !cfg.WasSetupRun || force {
				slog.Debug("Looks like this setup was never run, running setup now")
				fmt.Println("setting up package")
				err := tree.RunSetup(ctx)
				if err != nil {
					return errors.Wrap(err, "failed to run package setup")
				}
				cfg.WasSetupRun = true
				cfg.Save()
			} else {
				slog.Warn("Set up was run and I don't need to re-run")
			}

			return nil
		},
	}

	root.AddCommand(SetupCommand)
	SetupCommand.Flags().BoolVarP(&force, "force", "f", false, "Force setup to run, even if not needed")
}
