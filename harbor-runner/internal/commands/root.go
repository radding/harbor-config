package commands

import (
	"os"

	"github.com/radding/harbor-runner/internal/executor"
	"github.com/radding/harbor-runner/internal/telemetry"
	"github.com/spf13/cobra"
)

var machineReadableLogs *bool
var logLevel *telemetry.LogLevel = telemetry.LogLevelPtr(telemetry.InfoLevel)

var rootCmd = &cobra.Command{
	Short: "Harbor is a tool to manage workspaces for projects",
	Long:  `Harbor is a workspace management and build tool that enables developers to manage their projects more effectively.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

type RootExecutor struct {
	Exec executor.Executor
}

func (r *RootExecutor) Initialize() error {
	machineReadableLogs = rootCmd.PersistentFlags().BoolP("machine-readable", "m", false, "Produce machine readable JSON logs?")
	rootCmd.PersistentFlags().VarP(logLevel, "log-level", "v", "The Log level to set the logger to. Can be: Fatal, Error, Warn, Info, Debug, Http, Metrics, and Trace")
	rootCmd.ParseFlags(os.Args)
	createRunCommand(rootCmd, r.Exec)
	createSetupCommand(rootCmd, r.Exec)
	telemetry.ConfigureLogs(*machineReadableLogs, *logLevel)
	return nil
}

func (r *RootExecutor) Execute() error {
	return rootCmd.Execute()
}
