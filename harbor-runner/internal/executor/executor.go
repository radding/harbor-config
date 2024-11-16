package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/radding/harbor-runner/internal/application"
	"github.com/radding/harbor-runner/internal/taskgraph"
	"github.com/radding/harbor-runner/internal/telemetry"
)

type ExecutionResponse struct {
	WasCached bool
	Artifacts []struct {
		Name     string
		Location string
	}
	Error error
}

type ExecutionRequest struct {
	Kind          string
	TaskID        string
	WithCache     bool
	ForceClean    bool
	CacheLocation string
	WorkingDir    string
	WorkspaceRoot string
	Options       json.RawMessage
}

type ExecutionElement interface {
	Execute(ctx context.Context, msg ExecutionRequest) (ExecutionResponse, error)
}

type executor struct {
	executors map[string]ExecutionElement
}

// Initialize implements Executor.

type ExecutorOptions struct {
	executors map[string]ExecutionElement
}

func (e *executor) Register(kind string, exec any) {
	realExec := exec.(ExecutionElement)
	e.executors[kind] = realExec
}

type ExecutionOption func(e *ExecutorOptions) *ExecutorOptions

type Executor interface {
	taskgraph.Executor
	application.Initializer
}

func New(opts ...ExecutionOption) Executor {
	realOpt := &ExecutorOptions{
		executors: map[string]ExecutionElement{},
	}
	for _, opt := range opts {
		realOpt = opt(realOpt)
	}
	e := &executor{
		executors: realOpt.executors,
	}

	exec := &ExecCommand{}
	exec.RegsiterWith(e)

	return e
}

func WithKind(kind string, exec ExecutionElement) ExecutionOption {
	return func(e *ExecutorOptions) *ExecutorOptions {
		e.executors[kind] = exec
		return e
	}
}

func (e *executor) Execute(ctx context.Context, kind string, opts json.RawMessage) error {
	telemetry.Trace("running an executor", slog.String("kind", kind))
	withCache, ok := ctx.Value("WithCache").(bool)
	if !ok {
		withCache = true
	}
	forceClean, ok := ctx.Value("ForceClean").(bool)
	if !ok {
		forceClean = false
	}
	cacheLocation, ok := ctx.Value("CacheLocation").(string)
	if !ok {
		return fmt.Errorf("Failed to execute %s. Cache location not in context", kind)
	}
	workingDir, ok := ctx.Value("WorkingDir").(string)
	if !ok {
		return fmt.Errorf("Failed to execute %s. Working location not in context", kind)
	}
	taskID, ok := ctx.Value("TaskID").(string)
	if !ok {
		return fmt.Errorf("Failed to execute %s. Task ID not in context", kind)
	}
	workspaceRoot, ok := ctx.Value("workspaceRoot").(string)
	if !ok {
		slog.Warn("not in a workspace")
	}
	executor, ok := e.executors[kind]
	if !ok {
		return fmt.Errorf("No executor for kind %s", kind)
	}
	resp, err := executor.Execute(ctx, ExecutionRequest{
		Kind:          kind,
		WithCache:     withCache,
		ForceClean:    forceClean,
		CacheLocation: cacheLocation,
		WorkingDir:    workingDir,
		WorkspaceRoot: workspaceRoot,
		Options:       opts,
		TaskID:        taskID,
	})
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return errors.Wrap(err, "failed to execute")
	}

	if len(resp.Artifacts) > 0 {
		slog.Debug("Here I will move artifacts")
	}
	return nil
}
