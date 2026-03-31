package kk_scheduler

import (
	"context"
	"fmt"
	"log/slog"

	"gitee.com/cruvie/kk_go_kit/kk_id"
)

type StepCtl struct {
	addLog   func(message string)
	hasError bool
}

// Log adds a log message to the execution record
func (c *StepCtl) Log(message string) {
	c.addLog(message)
}

// Step represents a single execution step
type Step struct {
	name     string
	handler  func(ctl *StepCtl) error
	fallback func(ctl *StepCtl) error
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
func (t *TaskExecutor) AddStep(
	name string,
	handler func(ctl *StepCtl) error,
	fallback func(ctl *StepCtl) error,
) {
	t.steps = append(t.steps, &Step{
		name:     name,
		handler:  handler,
		fallback: fallback,
	})
}

// Run executes all steps sequentially
func (t *TaskExecutor) Run(ctx context.Context) error {
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
		input.SetId(t.id)
		input.SetJobId(t.jobId)
		_, err := t.client.TaskCreate(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to create execution record: %w", err)
		}
	}

	// Create shared StepCtl for all steps
	ctl := &StepCtl{
		addLog: func(message string) {
			in := &TaskAppendLog_Input{}
			in.SetId(t.id)
			in.SetLog(message)
			_, err := t.client.TaskAppendLog(ctx, in)
			if err != nil {
				slog.Error("failed to report log", "err", err)
			}
		},
	}
	ctl.Log("[🚀Task Start]")
	// Execute steps
	for _, step := range t.steps {
		ctl.Log(fmt.Sprintf("[🧩Step Begin %s]", step.name))
		err := step.handler(ctl)
		if err != nil {
			ctl.Log(fmt.Sprintf("❌%s", err.Error()))
			ctl.hasError = true
			if step.fallback != nil {
				ctl.Log(fmt.Sprintf("[🧯Step Fallback %s]", step.name))
				err := step.fallback(ctl)
				if err != nil {
					ctl.Log(fmt.Sprintf("❌%s", err.Error()))
				}
			}
			break
		}
	}

	in := &TaskUpdateStatus_Input{}
	in.SetId(t.id)
	if ctl.hasError {
		ctl.Log("[😭Task Failed]")
		in.SetStatus(TaskExecutionStatus_TASK_EXECUTION_STATUS_FAILED)
	} else {
		ctl.Log("[✅Task Finished]")
		in.SetStatus(TaskExecutionStatus_TASK_EXECUTION_STATUS_COMPLETED)
	}
	_, err := t.client.TaskUpdateStatus(ctx, in)
	if err != nil {
		slog.Error("failed to update task status", "err", err)
	}

	return nil
}
