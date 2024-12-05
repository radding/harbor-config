package builtins

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/radding/harbor-runner/internal/executor"
	packageconfig "github.com/radding/harbor-runner/internal/package-config"
	"github.com/radding/harbor-runner/internal/taskgraph"
	"github.com/radding/harbor-runner/internal/telemetry"
)

type localDependencyOptions struct {
	Path string `json:"path"`
}

type localDependency struct {
	taskGraph *taskgraph.ExecutionTree
	config    *packageconfig.Config
}

type LocalDependencyManager struct {
	locals map[string]*localDependency
}

func (l *LocalDependencyManager) Execute(ctx context.Context, msg executor.ExecutionRequest) (executor.ExecutionResponse, error) {
	telemetry.Trace("Encountered a local dependency, will make this work")
	opts := &localDependencyOptions{}
	err := json.Unmarshal(msg.Options, opts)
	if err != nil {
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to unmarshal JSON")
	}
	wd, err := os.Getwd()
	if err != nil {
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to get working directory")
	}
	pth := path.Join(wd, opts.Path, "./.harborrc.ts")
	slog.Debug("attempting to load config", slog.String("locaation", pth))
	conf, err := packageconfig.LoadConfig(pth)
	if err != nil {
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to load config for local dependency")
	}
	tree, err := taskgraph.CreateTreeFromConfig(&conf, msg.Task.GetExecutor())
	if err != nil {
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to get task graph for local dep")
	}
	l.locals[opts.Path] = &localDependency{
		taskGraph: tree,
		config:    &conf,
	}
	return executor.ExecutionResponse{}, nil
}

func (n *LocalDependencyManager) RegisterWith(reg executor.Registery) {
	reg.Register("harbor.dev/LocalDependency", n)
}
