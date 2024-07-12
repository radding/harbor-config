package telemetry

import (
	"time"

	"github.com/rs/zerolog/log"
)

func TimeWithError(opName string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start)
	log.Trace().
		Int64("duration", duration.Milliseconds()).
		Str("operation", opName).
		Msg("Telemetry Data")
	return err
}

func Time(opName string, fn func()) {
	TimeWithError(opName, func() error {
		fn()
		return nil
	})
}
