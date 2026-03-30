package task_executor

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/cruvie/kk-scheduler/internal/models"
	"github.com/cruvie/kk-scheduler/internal/store_driver"
	"gorm.io/gorm"
)

// Step represents a single execution step
type Step struct {
	name    string
	index   int
	handler func(ctl *StepCtl)
}

// TaskExecutor manages task lifecycle and step execution
type TaskExecutor struct {
	name              string
	steps             []*Step
	timeout           time.Duration
	grpcClientBuilder GRPCClientBuilder
	store             StoreDriver
	mu                sync.Mutex

	// runtime state
	id           string
	ctx          context.Context
	cancelSignal context.CancelFunc
	status       models.TaskExecutionStatus
	startTime    time.Time
	stopped      bool // marks if stop was called
}

// NewTaskStep creates a new task executor
func NewTaskStep(name string, db *gorm.DB, opts ...TaskOption) *TaskExecutor {
	t := &TaskExecutor{
		name:   name,
		status: models.TaskExecutionStatusRunning,
	}

	// Set default store if db is provided
	if db != nil {
		t.store = store_driver.NewPostgresStore(db)
	}

	// Apply options
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
	t.startTime = time.Now()

	// Ensure store is set
	if t.store == nil {
		return fmt.Errorf("store not configured, use WithStore or provide db to NewTaskStep")
	}

	// Create execution record
	id, err := t.store.Create(t.name, models.TaskExecutionStatusRunning, t.startTime.Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("failed to create execution record: %w", err)
	}
	t.id = id

	// Setup timeout context
	if t.timeout > 0 {
		t.ctx, t.cancelSignal = context.WithTimeout(context.Background(), t.timeout)
	} else {
		t.ctx, t.cancelSignal = context.WithCancel(context.Background())
	}

	// Create global StepCtl
	ctl := &StepCtl{
		addLog: func(message string) {
			formatted := fmt.Sprintf("[步骤%d: %s] %s\n", 0, "", message)
			t.store.AppendLog(t.id, formatted)
		},
		stop: func() {
			t.mu.Lock()
			t.stopped = true
			t.mu.Unlock()
			t.cancelSignal()
		},
	}

	// Execute steps
	for _, step := range t.steps {
		// Pre-step termination check
		if t.shouldStop(step.index) {
			t.finish(t.status)
			return t.statusToError()
		}

		// Update ctl for this step
		ctl.addLog = func(message string) {
			formatted := fmt.Sprintf("[步骤%d: %s] %s\n", step.index, step.name, message)
			t.store.AppendLog(t.id, formatted)
		}

		// Execute step with panic recovery
		if err := t.executeStep(step, ctl); err != nil {
			t.finish(models.TaskExecutionStatusFailed)
			return err
		}
	}

	t.finish(models.TaskExecutionStatusCompleted)
	return nil
}

// shouldStop checks if task should stop before executing next step
func (t *TaskExecutor) shouldStop(stepIndex int) bool {
	// Check timeout first
	if t.ctx.Err() == context.DeadlineExceeded {
		t.status = models.TaskExecutionStatusTimeout
		return true
	}

	// Check manual stop signal
	t.mu.Lock()
	stopped := t.stopped
	t.mu.Unlock()
	if stopped {
		t.status = models.TaskExecutionStatusCancelled
		return true
	}

	// Check gRPC termination instruction (if configured)
	if t.grpcClientBuilder != nil {
		if t.checkTerminationInstruction() {
			t.status = models.TaskExecutionStatusCancelled
			return true
		}
	}

	return false
}

// checkTerminationInstruction queries kk-scheduler for stop signal
func (t *TaskExecutor) checkTerminationInstruction() bool {
	client, err := t.grpcClientBuilder()
	if err != nil {
		slog.Warn("failed to get gRPC client for termination check", "err", err)
		return false // Don't block execution on gRPC failure
	}

	// Note: The actual RPC for checking termination would be defined in proto
	// For now, this is a placeholder that would call a CheckTermination RPC
	// The client is created but not used since the RPC doesn't exist yet
	slog.Debug("gRPC client created for termination check", "client", client)
	return false
}

// executeStep runs a single step with panic recovery
func (t *TaskExecutor) executeStep(step *Step, ctl *StepCtl) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("step panicked", "step", step.name, "panic", r)
			ctl.addLog(fmt.Sprintf("PANIC: %v", r))
			t.finish(models.TaskExecutionStatusFailed)
		}
	}()

	step.handler(ctl)
	return nil
}

// finish updates the final status and timestamp
func (t *TaskExecutor) finish(status models.TaskExecutionStatus) {
	t.status = status
	finishedAt := time.Now().Format(time.RFC3339Nano)
	t.store.UpdateStatus(t.id, status, finishedAt)
}

// statusToError converts status to appropriate error type
func (t *TaskExecutor) statusToError() error {
	switch t.status {
	case models.TaskExecutionStatusTimeout:
		return ErrTaskTimeout
	case models.TaskExecutionStatusCancelled:
		return ErrTaskCancelled
	default:
		return fmt.Errorf("task ended with status: %s", t.status)
	}
}
