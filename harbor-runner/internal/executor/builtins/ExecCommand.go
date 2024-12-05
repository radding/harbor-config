package builtins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
	"github.com/radding/harbor-runner/internal/executor"
	packageconfig "github.com/radding/harbor-runner/internal/package-config"
)

type ExecCommand struct {
}

// RegisterWith implements executor.ExecutionElement.
func (e *ExecCommand) RegisterWith(reg executor.Registery) {
	reg.Register("harbor.dev/ExecCommand", e)
}

type ExecOptions struct {
	Executable string            `json:"executable"`
	Args       []string          `json:"args"`
	Inputs     []string          `json:"inputs"`
	Env        map[string]string `json:"env"`
}

func (e *ExecCommand) Execute(ctx context.Context, msg executor.ExecutionRequest) (executor.ExecutionResponse, error) {
	done := make(chan executor.ExecutionResponse)
	slog.Debug("starting task", slog.String("working_dir", msg.WorkingDir), slog.String("name", msg.Task.ID))
	errChan := make(chan error)
	defer close(errChan)
	defer close(done)

	opts := ExecOptions{}
	err := json.Unmarshal(msg.Options, &opts)
	if err != nil {
		slog.Error("failed to parse options", slog.String("error", err.Error()))
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to parse options JSON")
	}

	taskName := msg.Task.ID

	infoCached, err := msg.Cache.Get("info.log", os.Stdout)
	if err != nil {
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to get from cache")
	}
	errorCached, err := msg.Cache.Get("error.log", os.Stderr)
	if err != nil {
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to get error log from cache")
	}
	if infoCached || errorCached {
		slog.Info("replayed from cache", slog.String("task_name", taskName))
		return executor.ExecutionResponse{}, nil
	}

	infoBuff := new(bytes.Buffer)
	errorBuff := new(bytes.Buffer)
	cmd := exec.Command(opts.Executable, opts.Args...)
	env := os.Environ()
	for key, val := range opts.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, val))
	}
	cmd.Env = env

	cmd.Stdout = &PipedLogger{
		logger: packageconfig.NewPipedLogger(slog.Info, slog.String("task_name", taskName)),
		fi:     infoBuff,
	}
	cmd.Stderr = &PipedLogger{
		logger: packageconfig.NewPipedLogger(slog.Error, slog.String("task_name", taskName)),
		fi:     errorBuff,
	}
	err = cmd.Start()
	if err != nil {
		slog.Error("failed to start command", slog.String("component", "harbor.dev/ExecCommand"), slog.String("error", err.Error()), slog.String("command", opts.Executable))
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to start command")
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			slog.Error("failed to run command", slog.String("component", "harbor.dev/ExecCommand"), slog.String("error", err.Error()), slog.String("command", opts.Executable))
			errChan <- err
			return
		}
		done <- executor.ExecutionResponse{
			WasCached: true,
			Artifacts: []struct {
				Name     string
				Location string
			}{},
			Error: nil,
		}
	}()
	select {
	case resp := <-done:
		msg.Cache.Add("info.log", infoBuff)
		msg.Cache.Add("error.log", errorBuff)
		return resp, nil
	case err := <-errChan:
		return executor.ExecutionResponse{}, errors.Wrap(err, "failed to execute command")
	case <-ctx.Done():
		cmd.Process.Signal(syscall.SIGINT)
		return executor.ExecutionResponse{}, fmt.Errorf("%s was canceled", msg.Kind)
	}
}

type PipedLogger struct {
	fi     io.Writer
	logger io.Writer
}

func (p *PipedLogger) Write(b []byte) (int, error) {
	n, err := p.logger.Write(b)
	if err != nil {
		return n, errors.Wrap(err, "failed to log message")
	}
	return p.fi.Write(b)
}
