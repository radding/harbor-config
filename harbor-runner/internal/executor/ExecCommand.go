package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/pkg/errors"
	packageconfig "github.com/radding/harbor-runner/internal/package-config"
)

type ExecCommand struct {
}

type ExecOptions struct {
	Executable string            `json:"executable"`
	Args       []string          `json:"args"`
	Inputs     []string          `json:"inputs"`
	Env        map[string]string `json:"env"`
}

type Register interface {
	Register(kind string, elem any)
}

func (e *ExecCommand) RegsiterWith(r Register) {
	r.Register("harbor.dev/ExecCommand", e)
}

func (e *ExecCommand) Execute(ctx context.Context, msg ExecutionRequest) (ExecutionResponse, error) {
	done := make(chan ExecutionResponse)
	errChan := make(chan error)
	defer close(errChan)
	defer close(done)

	opts := ExecOptions{}
	err := json.Unmarshal(msg.Options, &opts)
	if err != nil {
		slog.Error("failed to parse options", slog.String("error", err.Error()))
		return ExecutionResponse{}, errors.Wrap(err, "failed to parse options JSON")
	}
	cacheLocation := ""
	taskName, ok := ctx.Value("task_name").(string)
	if !ok {
		taskName = msg.Kind
	}
	if len(opts.Inputs) == 0 {
		cacheLocation = msg.CacheLocation
	}
	pth := path.Join(cacheLocation, msg.TaskID)
	slog.Debug("attempting to execute", slog.String("cachePath", pth))
	mkcmd := exec.Command("mkdir", "-p", pth)
	err = mkcmd.Run()
	if err != nil {
		return ExecutionResponse{}, errors.Wrap(err, "failed to make all dirs")
	}
	infoLogCachePath := path.Join(cacheLocation, msg.TaskID, "info.log")
	errorLogCachePath := path.Join(cacheLocation, msg.TaskID, "error.log")
	if data, err := os.ReadFile(infoLogCachePath); err == nil {
		slog.Info("replaying element from cache", slog.String("task_name", taskName))
		slog.Info(string(data), slog.String("task_name", taskName))
		d2, err := os.ReadFile(errorLogCachePath)
		if err != nil && !errors.Is(os.ErrNotExist, err) {
			return ExecutionResponse{}, errors.Wrap(err, "failed to open error file")
		}
		if len(d2) > 0 {
			slog.Error(string(d2), slog.String("task_name", taskName))
		}
		return ExecutionResponse{}, nil
	} else if errors.Is(os.ErrNotExist, err) {
		err = nil
	}

	infoFi, err := os.OpenFile(infoLogCachePath, os.O_CREATE, 0666)
	if err != nil {
		return ExecutionResponse{}, errors.Wrap(err, "failed to open info cache")
	}
	defer infoFi.Close()
	errorFi, err := os.OpenFile(errorLogCachePath, os.O_CREATE, 0666)
	if err != nil {
		return ExecutionResponse{}, errors.Wrap(err, "failed to open error cache")
	}
	defer errorFi.Close()
	cmd := exec.Command(opts.Executable, opts.Args...)
	env := os.Environ()
	for key, val := range opts.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, val))
	}
	cmd.Env = env

	cmd.Stdout = &PipedLogger{
		logger: packageconfig.NewPipedLogger(slog.Info, slog.String("task_name", taskName)),
		fi:     infoFi,
	}
	cmd.Stderr = &PipedLogger{
		logger: packageconfig.NewPipedLogger(slog.Error, slog.String("task_name", taskName)),
		fi:     errorFi,
	}
	err = cmd.Start()
	if err != nil {
		slog.Error("failed to start command", slog.String("component", "harbor.dev/ExecCommand"), slog.String("error", err.Error()), slog.String("command", opts.Executable))
		return ExecutionResponse{}, errors.Wrap(err, "failed to start command")
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			slog.Error("failed to run command", slog.String("component", "harbor.dev/ExecCommand"), slog.String("error", err.Error()), slog.String("command", opts.Executable))
			errChan <- err
			return
		}
		done <- ExecutionResponse{
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
		return resp, nil
	case err := <-errChan:
		cmd.Process.Signal(syscall.SIGINT)
		return ExecutionResponse{}, errors.Wrap(err, "failed to execute command")
	case <-ctx.Done():
		cmd.Process.Signal(syscall.SIGINT)
		return ExecutionResponse{}, fmt.Errorf("%s was canceled", msg.Kind)
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
