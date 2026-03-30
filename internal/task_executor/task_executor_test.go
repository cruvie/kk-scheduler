package task_executor

import (
	"testing"
	"time"

	"github.com/cruvie/kk-scheduler/internal/models"
	"github.com/cruvie/kk-scheduler/internal/store_driver"
	"github.com/stretchr/testify/assert"
)

func TestNewTaskStep(t *testing.T) {
	mockStore := store_driver.NewMockStore()
	task := NewTaskStep("test-task", nil, WithStore(mockStore))

	assert.NotNil(t, task)
	assert.Equal(t, "test-task", task.name)
	assert.Equal(t, mockStore, task.store)
}

func TestAddStep(t *testing.T) {
	task := NewTaskStep("test", nil, WithStore(store_driver.NewMockStore()))

	task.AddStep("step1", func(ctl *StepCtl) {})
	task.AddStep("step2", func(ctl *StepCtl) {})

	assert.Len(t, task.steps, 2)
	assert.Equal(t, "step1", task.steps[0].name)
	assert.Equal(t, 0, task.steps[0].index)
	assert.Equal(t, "step2", task.steps[1].name)
	assert.Equal(t, 1, task.steps[1].index)
}

func TestStepCtlAddLog(t *testing.T) {
	mockStore := store_driver.NewMockStore()
	task := NewTaskStep("test", nil, WithStore(mockStore))

	// Create a record first (simulating what Run() does)
	mockStore.TaskCreate("test", models.TaskExecutionStatusRunning)

	ctl := &StepCtl{
		addLog: func(msg string) { mockStore.TaskAppendLog(task.name, msg) },
		stop:   func() {},
	}

	ctl.AddLog("test message")

	logs := mockStore.GetLogs("test")
	assert.Len(t, logs, 1)
	assert.Contains(t, logs[0], "test message")
}

func TestStepCtlStop(t *testing.T) {
	stopCalled := false
	ctl := &StepCtl{
		addLog: func(msg string) {},
		stop:   func() { stopCalled = true },
	}

	ctl.Stop()
	assert.True(t, stopCalled)
}

func TestStepCtlNilFunctions(t *testing.T) {
	ctl := &StepCtl{}

	// Should not panic when functions are nil
	ctl.AddLog("test")
	ctl.Stop()
}

func TestRunBasic(t *testing.T) {
	mockStore := store_driver.NewMockStore()
	task := NewTaskStep("basic-test", nil, WithStore(mockStore))

	task.AddStep("step1", func(ctl *StepCtl) {
		ctl.AddLog("executing step1")
	})
	task.AddStep("step2", func(ctl *StepCtl) {
		ctl.AddLog("executing step2")
	})

	err := task.Run()
	assert.NoError(t, err)

	// Check execution record
	assert.Equal(t, 1, mockStore.CreatedCount)
	record := mockStore.GetRecord("basic-test")
	assert.Equal(t, models.TaskExecutionStatusCompleted, record.Status)

	// Check logs
	logs := mockStore.GetLogs("basic-test")
	assert.Len(t, logs, 2)
}

func TestRunEmptySteps(t *testing.T) {
	mockStore := store_driver.NewMockStore()
	task := NewTaskStep("empty-test", nil, WithStore(mockStore))

	err := task.Run()
	assert.NoError(t, err)

	record := mockStore.GetRecord("empty-test")
	assert.Equal(t, models.TaskExecutionStatusCompleted, record.Status)
}

func TestRunWithTimeout(t *testing.T) {
	mockStore := store_driver.NewMockStore()
	task := NewTaskStep("timeout-test", nil,
		WithStore(mockStore),
		WithTimeout(50*time.Millisecond),
	)

	// Step 1 executes quickly
	task.AddStep("quick", func(ctl *StepCtl) {
		ctl.AddLog("quick step")
	})

	// Sleep to trigger timeout before step 2
	task.AddStep("slow", func(ctl *StepCtl) {
		time.Sleep(100 * time.Millisecond)
		ctl.AddLog("slow step")
	})

	task.AddStep("after-slow", func(ctl *StepCtl) {
		ctl.AddLog("this should not execute")
	})

	err := task.Run()
	// Either timeout or cancelled depending on timing
	assert.Error(t, err)
}

func TestRunPanicRecovery(t *testing.T) {
	mockStore := store_driver.NewMockStore()
	task := NewTaskStep("panic-test", nil, WithStore(mockStore))

	task.AddStep("panic-step", func(ctl *StepCtl) {
		ctl.AddLog("about to panic")
		panic("test panic")
	})

	err := task.Run()
	// Panic recovery returns no error from Run(), but status is failed
	assert.NoError(t, err)

	record := mockStore.GetRecord("panic-test")
	assert.Equal(t, models.TaskExecutionStatusFailed, record.Status)

	logs := mockStore.GetLogs("panic-test")
	assert.True(t, len(logs) >= 1)
}

func TestRunManualStop(t *testing.T) {
	mockStore := store_driver.NewMockStore()
	task := NewTaskStep("stop-test", nil, WithStore(mockStore))

	task.AddStep("step1", func(ctl *StepCtl) {
		ctl.AddLog("step1")
		ctl.Stop() // Manually stop
	})

	task.AddStep("step2", func(ctl *StepCtl) {
		ctl.AddLog("this should not execute")
	})

	err := task.Run()
	assert.ErrorIs(t, err, ErrTaskCancelled)

	record := mockStore.GetRecord("stop-test")
	assert.Equal(t, models.TaskExecutionStatusCancelled, record.Status)

	logs := mockStore.GetLogs("stop-test")
	assert.Len(t, logs, 1) // Only step1 log
}

func TestNoStoreError(t *testing.T) {
	task := NewTaskStep("no-store-test", nil)

	task.AddStep("step1", func(ctl *StepCtl) {})

	err := task.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "store not configured")
}

func TestWithTimeoutOption(t *testing.T) {
	task := NewTaskStep("test", nil, WithTimeout(10*time.Second))

	assert.Equal(t, 10*time.Second, task.timeout)
}

func TestWithStoreOption(t *testing.T) {
	mockStore := store_driver.NewMockStore()
	task := NewTaskStep("test", nil, WithStore(mockStore))

	assert.Equal(t, mockStore, task.store)
}

func TestMultipleStepsWithLogs(t *testing.T) {
	mockStore := store_driver.NewMockStore()
	task := NewTaskStep("multi-log-test", nil, WithStore(mockStore))

	task.AddStep("initialize", func(ctl *StepCtl) {
		ctl.AddLog("initializing")
		ctl.AddLog("initialized")
	})

	task.AddStep("process", func(ctl *StepCtl) {
		ctl.AddLog("processing data")
	})

	task.AddStep("finalize", func(ctl *StepCtl) {
		ctl.AddLog("finalizing")
	})

	err := task.Run()
	assert.NoError(t, err)

	logs := mockStore.GetLogs("multi-log-test")
	assert.Len(t, logs, 4)

	// Check log format contains step info
	assert.Contains(t, logs[0], "[步骤0: initialize]")
	assert.Contains(t, logs[2], "[步骤1: process]")
	assert.Contains(t, logs[3], "[步骤2: finalize]")
}
