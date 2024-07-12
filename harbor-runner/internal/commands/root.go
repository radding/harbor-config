package commands

import (
	"os"

	"github.com/radding/harbor-runner/internal/telemetry"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var machineReadableLogs *bool
var logLevel *telemetry.LogLevel = telemetry.LogLevelPtr(telemetry.LogLevel(zerolog.InfoLevel))

func init() {
	machineReadableLogs = rootCmd.PersistentFlags().BoolP("machine-readable", "m", false, "Produce machine readable JSON logs?")
	rootCmd.PersistentFlags().VarP(logLevel, "log-level", "v", "The Log level to set the logger to. Can be: Panic, Fatal, Error, Warn, Info, Debug, and Trace")
}

var rootCmd = &cobra.Command{
	Short: "Harbor is a tool to manage workspaces for projects",
	Long:  `Harbor is a workspace management and build tool that enables developers to manage their projects more effectively.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

type RootExecutor struct{}

func (r *RootExecutor) Initialize() error {
	rootCmd.ParseFlags(os.Args)
	telemetry.ConfigureLogs(*machineReadableLogs, *logLevel)
	return nil
}

func (r *RootExecutor) Execute() error {
	return rootCmd.Execute()
}
