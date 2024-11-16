package taskgraph

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockExecutor struct {
	executionOrder []string
	mockFunc       func(kind string) error
}

func (m *MockExecutor) Execute(ctx context.Context, kind string, opts json.RawMessage) error {
	m.executionOrder = append(m.executionOrder, kind)
	if m.mockFunc != nil {
		return m.mockFunc(kind)
	}
	return nil
}

func TestExecutionOrder(t *testing.T) {
	assert := assert.New(t)
	executor := &MockExecutor{}
	task1 := &Task{
		executor: executor,
		Kind:     "test1",
	}
	task2 := &Task{
		executor: executor,
		Kind:     "test2",
	}
	task3 := &Task{
		executor:     executor,
		Kind:         "test3",
		Dependencies: []*Task{task1, task2},
	}

	rootTask := &Task{
		executor:     executor,
		Kind:         "test4",
		Dependencies: []*Task{task3, task1},
	}
	err := rootTask.Execute(context.Background())
	assert.NoError(err)
	assert.Equal([]string{"test1", "test2", "test3", "test4"}, executor.executionOrder)
}

func TestErrorCancelsEverything(t *testing.T) {
	assert := assert.New(t)
	executor := &MockExecutor{
		mockFunc: func(kind string) error {
			if kind == "blow_up" {
				return errors.New("We don't like your kind around here")
			}
			return nil
		},
	}
	task1 := &Task{
		executor: executor,
		Kind:     "test1",
	}
	task2 := &Task{
		executor: executor,
		Kind:     "test2",
	}

	exploder := &Task{
		executor:     executor,
		Kind:         "blow_up",
		Dependencies: []*Task{task1, task2},
	}
	task3 := &Task{
		executor:     executor,
		Kind:         "test3",
		Dependencies: []*Task{exploder},
	}

	rootTask := &Task{
		executor:     executor,
		Kind:         "test4",
		Dependencies: []*Task{task3, task1},
	}
	err := rootTask.Execute(context.Background())
	assert.Error(err)
	assert.Equal([]string{"test1", "test2", "blow_up"}, executor.executionOrder)
}
