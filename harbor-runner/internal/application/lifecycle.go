package application

import (
	"errors"
	"fmt"
	"log/slog"
)

type Initializer interface {
	Initialize() error
}

type Executor interface {
	Execute() error
}

type Cleaner interface {
	Clean() error
}

type Application struct {
	initializers []Initializer
	executor     Executor
	cleaners     []Cleaner
}

func (a *Application) Initialize() error {
	errs := []error{}
	for _, init := range a.initializers {
		err := init.Initialize()
		if err != nil {
			slog.Warn(fmt.Sprintf("error initializing: %s", err))
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (a *Application) Exectute() error {
	if a.executor == nil {
		return fmt.Errorf("can't execute program because no executor defined")
	}
	return a.executor.Execute()
}

func (a *Application) Clean() error {
	errs := []error{}
	for _, init := range a.cleaners {
		err := init.Clean()
		if err != nil {
			slog.Warn(fmt.Sprintf("error cleaning up: %s", err))
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (a *Application) Register(m any) {
	foundAtleastOne := false
	if mod, ok := m.(Cleaner); ok {
		a.cleaners = append(a.cleaners, mod)
		foundAtleastOne = true
	}
	if mod, ok := m.(Initializer); ok {
		a.initializers = append(a.initializers, mod)
		foundAtleastOne = true
	}
	if mod, ok := m.(Executor); ok {
		if a.executor != nil {
			slog.Warn("Attempting to re-register an executor, I will not do this!")
			return
		}
		a.executor = mod
		foundAtleastOne = true
	}
	if !foundAtleastOne {
		panic("attempted to register an unknown interface")
	}
}

func New() *Application {
	return NewWithExecutor(nil)
}

func NewWithExecutor(exec Executor) *Application {
	return &Application{
		executor: exec,
	}
}
