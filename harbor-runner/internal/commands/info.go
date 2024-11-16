package commands

import (
	"bytes"
	"log/slog"
	"text/template"

	_ "embed"

	"github.com/pkg/errors"
	packageconfig "github.com/radding/harbor-runner/internal/package-config"
	"github.com/spf13/cobra"
)

var output string

func init() {
	rootCmd.AddCommand(InfoCommand)
	InfoCommand.Flags().StringVarP(&output, "output", "o", "text", "choose output for which the info")
}

//go:embed templates/info.tmpl
var infoTmpl string

var InfoCommand = &cobra.Command{
	Use:   "info",
	Short: "Get information on the current workspace/project",
	Long: `Print out the information on this workspace/project. 
	This will give information on dependencies, commands, local cache file, version, stability etc.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := packageconfig.GetConfig()
		if cfg == nil {
			return errors.New("failed to run command, no configuration found")
		}
		if output == "text" {
			tmpl, err := template.New("textTempl").Parse(infoTmpl)
			if err != nil {
				return errors.Wrap(err, "failed to parse info template")
			}
			buf := new(bytes.Buffer)
			err = tmpl.Execute(buf, cfg)
			if err != nil {
				return errors.Wrap(err, "Failed to execute template")
			}
			slog.Info(buf.String())
		}
		return nil
	},
}
