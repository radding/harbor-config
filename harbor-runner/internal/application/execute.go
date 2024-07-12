package application

import (
	"errors"
	"fmt"

	"github.com/radding/harbor-runner/internal/telemetry"
)

func RunApplication(app *Application) (err error) {
	defer func() {
		newErr := telemetry.TimeWithError("cleanup", func() error {
			return app.Clean()
		})
		if err != nil && newErr != nil {
			err = errors.Join(newErr, fmt.Errorf("failed to clean up and has other errors: %w", err))
		} else if newErr != nil {
			err = fmt.Errorf("failed to clean up: %w", newErr)
		}
	}()

	err = telemetry.TimeWithError("initialization", func() error {
		return app.Initialize()
	})

	if err != nil {
		err = fmt.Errorf("failed to initialize: %w", err)
		return
	}

	err = telemetry.TimeWithError("execution", func() error {
		return app.Exectute()
	})
	return
}
