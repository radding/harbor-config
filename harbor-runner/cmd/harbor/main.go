package main

import (
	"os"

	"github.com/radding/harbor-runner/internal/application"
	"github.com/radding/harbor-runner/internal/cfg"
	"github.com/radding/harbor-runner/internal/commands"
	exec "github.com/radding/harbor-runner/internal/executor"
	"github.com/radding/harbor-runner/internal/executor/builtins"
	packageconfig "github.com/radding/harbor-runner/internal/package-config"
	"github.com/radding/harbor-runner/internal/telemetry"
)

func main() {
	taskExecutor := exec.New()
	executor := &commands.RootExecutor{
		Exec: taskExecutor,
	}
	config := &cfg.ConfigLifeCycle{}

	app := application.New()
	app.Register(executor)
	app.Register(config)
	app.Register(&packageconfig.Lifecycle{})
	app.Register(builtins.New(taskExecutor))
	app.Register(taskExecutor)

	err := application.RunApplication(app)

	if err != nil {
		telemetry.Fatal("Unrecoverable error", err)
		os.Exit(1)
	}
}
