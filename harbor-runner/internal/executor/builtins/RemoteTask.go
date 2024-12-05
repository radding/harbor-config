package builtins

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/radding/harbor-runner/internal/executor"
)

type RemoteExecutor struct {
	localDeps *LocalDependencyManager
}

type remoteExecutorOptions struct {
	Dependency localDependencyOptions `json:"dependency"`
	Run        string                 `json:"run"`
	IsDepLocal bool                   `json:"isDepenedencyLocal"`
	Artifacts  []string               `json:"artifacts"`
	Inputs     []string               `json:"string"`
}

func (l *RemoteExecutor) Execute(ctx context.Context, msg executor.ExecutionRequest) (executor.ExecutionResponse, error) {
	opts := &remoteExecutorOptions{}

	err := json.Unmarshal(msg.Options, opts)
	if err != nil {
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to get options for remote task")
	}

	if opts.IsDepLocal {
		dep, ok := l.localDeps.locals[opts.Dependency.Path]
		if !ok {
			return executor.ExecutionResponse{}, errors.Errorf("did not load local depenedency at %s", opts.Dependency.Path)
		}
		ctx = dep.config.ConfigureContext(ctx)
		err := dep.taskGraph.RunTask(ctx, opts.Run)
		if err != nil {
			return executor.ExecutionResponse{}, errors.Wrap(err, "failed to run task")
		}
		return executor.ExecutionResponse{}, nil
	}
	return executor.ExecutionResponse{}, errors.New("not implemented")
}

func (n *RemoteExecutor) RegisterWith(reg executor.Registery) {
	reg.Register("harbor.dev/RemoteTask", n)
}
