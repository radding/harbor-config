package telemetry

import (
	"log/slog"
	"time"
)

func TimeWithError(opName string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start)
	Metrics(
		opName,
		slog.Int64("duration", duration.Milliseconds()),
		slog.Bool("failed", err != nil),
	)
	// log.Trace().
	// 	Int64("duration", duration.Milliseconds()).
	// 	Str("operation", opName).
	// 	Msg("Telemetry Data")
	return err
}

func Time(opName string, fn func()) {
	TimeWithError(opName, func() error {
		fn()
		return nil
	})
}
