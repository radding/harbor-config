package taskgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/pkg/errors"
	"github.com/radding/harbor-runner/internal/telemetry"
)

type Executor interface {
	Execute(ctx context.Context, kind string, opts json.RawMessage) error
}

type Task struct {
	ID            string
	Kind          string
	Options       json.RawMessage
	Dependencies  []*Task
	dependencySet map[string]bool
	executor      Executor
	done          bool
	err           error
}

func (t *Task) addDependency(t2 *Task) {
	if _, ok := t.dependencySet[t2.ID]; ok {
		return
	}
	t.dependencySet[t2.ID] = true
	t.Dependencies = append(t.Dependencies, t2)
}

func (t *Task) Execute(ctx context.Context) error {
	return telemetry.TimeWithError(fmt.Sprintf("executing task %s", t.ID), func() error {
		slog.Debug("Executing task", slog.String("task_id", t.ID))
		if t.done {
			slog.Debug("Task has already been done durring this run, returning early")
			return t.err
		}
		ctx = context.WithValue(ctx, "TaskID", t.ID)
		ctx, cancel := context.WithCancelCause(ctx)
		err := t.executeChildren(ctx)
		if err != nil {
			cancel(err)
			slog.Warn("task's children failed to execute", slog.String("task_id", t.ID), slog.String("error", err.Error()))
			t.done = true

			t.err = errors.Wrap(err, "failed to execute children")
			return t.err
		}
		err = t.executor.Execute(ctx, t.Kind, t.Options)
		if err != nil {
			slog.Warn("task failed to execute", slog.String("task_id", t.ID), slog.String("error", err.Error()))
			cancel(err)
			t.done = true
			t.err = errors.Wrap(err, "failed to execute task")
			return t.err
		}
		t.done = true
		return nil
	})
}

func (t *Task) executeChildren(ctx context.Context) error {
	ctx, cancel := context.WithCancelCause(ctx)
	waitCh := make(chan struct{})
	go func() {
		wg := &sync.WaitGroup{}
		for _, dep := range t.Dependencies {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := dep.Execute(ctx)
				if err != nil {
					telemetry.Trace("child died", slog.String("child_id", dep.ID), slog.String("task_id", dep.ID), slog.String("error", err.Error()))
					cancel(err)
				}
			}()
		}
		wg.Wait()
		waitCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		telemetry.Trace("Damn, someone failed :(")
		return ctx.Err()
	case <-waitCh:
		close(waitCh)
		telemetry.Trace("finished executing all children")
		return nil
	}
}
