package kk_scheduler

import (
	"context"
	"fmt"
	"log/slog"
)

// Step represents a single execution step
type Step struct {
	name    string
	index   int
	handler func(ctl *StepCtl)
}

// TaskExecutor manages task lifecycle and step execution
type TaskExecutor struct {
	name   string
	steps  []*Step
	client KKScheduleTaskExecutionClient
	status TaskExecutionStatus
}

// NewTaskExecutor creates a new task executor
func NewTaskExecutor(name string, opts ...TaskExecutorOption) *TaskExecutor {
	t := &TaskExecutor{
		name:   name,
		status: TaskExecutionStatusRunning,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// AddStep adds a step to the task
func (t *TaskExecutor) AddStep(name string, handler func(ctl *StepCtl)) {
	t.steps = append(t.steps, &Step{
		name:    name,
		index:   len(t.steps),
		handler: handler,
	})
}

// Run executes all steps sequentially
func (t *TaskExecutor) Run() error {
	ctx := context.Background()

	// Create execution record
	input := &TaskCreate_Input{}
	input.SetTaskName(t.name)
	input.SetStatus(string(t.status))
	_, err := t.client.TaskCreate(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create execution record: %w", err)
	}

	// Execute steps
	for _, step := range t.steps {
		ctl := &StepCtl{
			addLog: func(message string) {
				formatted := fmt.Sprintf("[步骤%d: %s] %s", step.index, step.name, message)
				logInput := &TaskAppendLog_Input{}
				logInput.SetTaskName(t.name)
				logInput.SetLog(formatted)
				_, err := t.client.TaskAppendLog(ctx, logInput)
				if err != nil {
					slog.Warn("failed to append log", "err", err)
				}
			},
		}

		// Execute step with panic recovery
		if err := t.executeStep(step, ctl); err != nil {
			t.finish(ctx, TaskExecutionStatusFailed)
			return err
		}

		// Check if status was changed by panic recovery
		if t.status == TaskExecutionStatusFailed {
			t.finish(ctx, TaskExecutionStatusFailed)
			return nil
		}
	}

	t.finish(ctx, TaskExecutionStatusCompleted)
	return nil
}

// executeStep runs a single step with panic recovery
func (t *TaskExecutor) executeStep(step *Step, ctl *StepCtl) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("step panicked", "step", step.name, "panic", r)
			ctl.AddLog(fmt.Sprintf("PANIC: %v", r))
			t.status = TaskExecutionStatusFailed
		}
	}()

	step.handler(ctl)
	return nil
}

// finish updates the final status
func (t *TaskExecutor) finish(ctx context.Context, status TaskExecutionStatus) {
	t.status = status
	updateInput := &TaskUpdateStatus_Input{}
	updateInput.SetTaskName(t.name)
	updateInput.SetStatus(string(status))
	_, err := t.client.TaskUpdateStatus(ctx, updateInput)
	if err != nil {
		slog.Warn("failed to update status", "err", err)
	}
}
