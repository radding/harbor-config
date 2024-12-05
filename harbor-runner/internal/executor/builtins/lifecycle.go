package builtins

import (
	"github.com/radding/harbor-runner/internal/executor"
)

type Lifecycle struct{}

func New(ex executor.Executor) *Lifecycle {
	execCommand := &ExecCommand{}
	pkgSetup := &PackageSetup{}
	noop := &Noop{}
	localDeps := &LocalDependencyManager{
		locals: map[string]*localDependency{},
	}
	remote := &RemoteExecutor{
		localDeps: localDeps,
	}

	ex.Accept(execCommand)
	ex.Accept(pkgSetup)
	ex.Accept(noop)
	ex.Accept(localDeps)
	ex.Accept(remote)

	return &Lifecycle{}
}

func (l *Lifecycle) Initialize() error {
	return nil
}
