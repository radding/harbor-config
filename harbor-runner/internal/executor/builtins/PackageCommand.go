package builtins

import (
	"context"
	"log/slog"

	"github.com/radding/harbor-runner/internal/executor"
	"github.com/radding/harbor-runner/internal/telemetry"
)

type PackageSetup struct{}

// Execute implements executor.ExecutionElement.
func (p *PackageSetup) Execute(ctx context.Context, msg executor.ExecutionRequest) (executor.ExecutionResponse, error) {
	slog.Debug("I am in the package setup, but I forgot what I am supposed to do right now", slog.Any("msg", string(msg.Options)))
	return executor.ExecutionResponse{
		WasCached: false,
	}, nil
}

func (p *PackageSetup) RegisterWith(reg executor.Registery) {
	reg.Register("harbor.dev/PackageSetup", p)
}

type Noop struct{}

func (n *Noop) Execute(ctx context.Context, msg executor.ExecutionRequest) (executor.ExecutionResponse, error) {
	telemetry.Trace("Noop executor executed. This is, well, a noop.")
	return executor.ExecutionResponse{
		WasCached: false,
	}, nil
}

func (n *Noop) RegisterWith(reg executor.Registery) {
	reg.Register("harbor.dev/noop", n)
}
