package main

import (
	"os"

	"github.com/radding/harbor-runner/internal/application"
	"github.com/radding/harbor-runner/internal/cfg"
	"github.com/radding/harbor-runner/internal/commands"
	"github.com/rs/zerolog/log"
)

func main() {
	executor := &commands.RootExecutor{}
	config := &cfg.ConfigLifeCycle{}

	app := application.New()
	app.Register(executor)
	app.Register(config)

	err := application.RunApplication(app)

	if err != nil {
		log.Fatal().Err(err).Msg("Unrecoverable error")
		os.Exit(1)
	}
}
