package kk_scheduler

import (
	"context"
	"fmt"
	"log/slog"

	"gitee.com/cruvie/kk_go_kit/kk_id"
)

// Step represents a single execution step
type Step struct {
	name    string
	handler func(ctl *StepCtl)
}

// TaskExecutor manages task lifecycle and step execution
type TaskExecutor struct {
	id     string
	jobId  string
	steps  []*Step
	client KKScheduleClient
}

// NewTaskExecutor creates a new task executor
func NewTaskExecutor(opts ...TaskExecutorOption) *TaskExecutor {
	t := &TaskExecutor{
		id: kk_id.GenUUID7(),
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
		handler: handler,
	})
}

// Run executes all steps sequentially
func (t *TaskExecutor) Run() error {
	ctx := context.Background()
	{
		if t.client == nil {
			return fmt.Errorf("scheduler client is not set")
		}
		if t.jobId == "" {
			return fmt.Errorf("job id is not set")
		}
	}
	{
		// Create execution record
		input := &TaskCreate_Input{}
		input.SetJobId(t.jobId)
		_, err := t.client.TaskCreate(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to create execution record: %w", err)
		}
	}

	// Execute steps
	for _, step := range t.steps {
		ctl := &StepCtl{
			addLog: func(message string) {
				formatted := fmt.Sprintf("[步骤%d: %s] %s", step.index, step.name, message)
				logInput := &TaskAppendLog_Input{}
				logInput.SetId(t.id)
				logInput.SetLog(formatted)
				_, err := t.client.TaskAppendLog(ctx, logInput)
				if err != nil {
					slog.Warn("failed to append log", "err", err)
				}
			},
		}

		// Execute step with panic recovery
		if err := t.executeStep(step, ctl); err != nil {
			t.finish(ctx, TaskExecutionStatus_TASK_EXECUTION_STATUS_FAILED)
			return err
		}

		// Check if status was changed by panic recovery
		if t.status == TaskExecutionStatus_TASK_EXECUTION_STATUS_FAILED {
			t.finish(ctx, TaskExecutionStatus_TASK_EXECUTION_STATUS_FAILED)
			return nil
		}
	}

	t.finish(ctx, TaskExecutionStatus_TASK_EXECUTION_STATUS_COMPLETED)
	return nil
}

// executeStep runs a single step with panic recovery
func (t *TaskExecutor) executeStep(step *Step, ctl *StepCtl) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("step panicked", "step", step.name, "panic", r)
			ctl.AddLog(fmt.Sprintf("PANIC: %v", r))
			t.status = TaskExecutionStatus_TASK_EXECUTION_STATUS_FAILED
		}
	}()

	step.handler(ctl)
	return nil
}

// finish updates the final status
func (t *TaskExecutor) finish(ctx context.Context, status TaskExecutionStatus) {
	t.status = status
	updateInput := &TaskUpdateStatus_Input{}
	updateInput.SetId(t.id)
	updateInput.SetStatus(status)
	_, err := t.client.TaskUpdateStatus(ctx, updateInput)
	if err != nil {
		slog.Warn("failed to update status", "err", err)
	}
}
