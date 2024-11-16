package taskgraph

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	packageconfig "github.com/radding/harbor-runner/internal/package-config"
)

type ExecutionTree struct {
	setupTask *Task
	tasks     map[string]*Task
}

func (e *ExecutionTree) RunSetup(ctx context.Context) error {
	return e.setupTask.Execute(ctx)
}
func (e *ExecutionTree) RunTask(ctx context.Context, taskName string) error {
	task, ok := e.tasks[taskName]
	if !ok {
		return fmt.Errorf("can not find task with name %s", taskName)
	}
	return task.Execute(ctx)
}

func CreateTreeFromConfig(cfg *packageconfig.Config, executor Executor) (*ExecutionTree, error) {
	setUpTask := &Task{
		Kind:          "harbor.dev/noop",
		ID:            "setup-mock-task",
		Dependencies:  []*Task{},
		dependencySet: map[string]bool{},
	}
	executionTree := &ExecutionTree{
		setupTask: setUpTask,
		tasks:     map[string]*Task{},
	}
	constructs := map[string]*Task{}

	var createTask func(taskID string) (*Task, error)
	createTask = func(taskID string) (*Task, error) {
		if elem, ok := constructs[taskID]; ok {
			return elem, nil
		}
		taskDef, ok := cfg.Constructs[taskID]
		if !ok {
			return nil, fmt.Errorf("failed to find construct of id %s", taskID)
		}
		t := &Task{
			executor:      executor,
			ID:            taskID,
			Kind:          taskDef.Kind,
			Options:       taskDef.Options,
			Dependencies:  []*Task{},
			done:          false,
			err:           nil,
			dependencySet: map[string]bool{},
		}
		for _, childID := range taskDef.DependsOn {
			childTask, err := createTask(childID)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create child of %s", taskID)
			}
			t.addDependency(childTask)
		}
		constructs[taskID] = t
		return t, nil
	}

	// // Iterate over all constructs and create a task
	// // for taskID, taskDef := range cfg.Constructs {
	for _, taskID := range cfg.Setup {
		t, err := createTask(taskID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create setup task")
		}
		executionTree.setupTask.addDependency(t)
	}

	for name, taskID := range cfg.Tasks {
		t, err := createTask(taskID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create task %s", name)
		}
		executionTree.tasks[name] = t
	}

	return executionTree, nil
}
