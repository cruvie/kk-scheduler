package kk_scheduler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cruvie/kk-scheduler/kk_scheduler"
	"github.com/stretchr/testify/assert"
)

func TestTaskExecutor_Run(t *testing.T) {
	t.Run("returns error when client is nil", func(t *testing.T) {
		executor := kk_scheduler.NewTaskExecutor(
			kk_scheduler.WithJobId("test-job"),
		)
		err := executor.Run(context.Background())
		assert.EqualError(t, err, "scheduler client is not set")
	})

	t.Run("returns error when jobId is empty", func(t *testing.T) {
		mockClient := NewMockStoreClient()
		executor := kk_scheduler.NewTaskExecutor(
			kk_scheduler.WithSchedulerClient(mockClient),
		)
		err := executor.Run(context.Background())
		assert.EqualError(t, err, "job id is not set")
	})

	t.Run("executes steps and updates status to completed", func(t *testing.T) {
		mockClient := NewMockStoreClient()
		jobId := "test-job-success"

		executor := kk_scheduler.NewTaskExecutor(
			kk_scheduler.WithSchedulerClient(mockClient),
			kk_scheduler.WithJobId(jobId),
		)

		stepExecuted := false
		executor.AddStep("step1", func(ctl *kk_scheduler.StepCtl) error {
			stepExecuted = true
			ctl.Log("step1 completed")
			return nil
		})

		err := executor.Run(context.Background())
		assert.NoError(t, err)
		assert.True(t, stepExecuted)

		logs := mockClient.GetStore().GetLogs(jobId)
		assert.Contains(t, logs, "🚀Starting task")
		assert.Contains(t, logs, "[🔥Starting step step1]")
		assert.Contains(t, logs, "step1 completed")
		assert.Contains(t, logs, "[✅Task Finished]")
	})

	t.Run("executes steps and updates status to failed when step returns error", func(t *testing.T) {
		mockClient := NewMockStoreClient()
		jobId := "test-job-failed"

		executor := kk_scheduler.NewTaskExecutor(
			kk_scheduler.WithSchedulerClient(mockClient),
			kk_scheduler.WithJobId(jobId),
		)

		executor.AddStep("failing-step", func(ctl *kk_scheduler.StepCtl) error {
			return errors.New("step failed")
		})

		err := executor.Run(context.Background())
		assert.NoError(t, err)

		logs := mockClient.GetStore().GetLogs(jobId)
		assert.Contains(t, logs, "❌step failed")
		assert.Contains(t, logs, "[❌😭🤬Task Failed]")
	})

	t.Run("continues executing steps after failure", func(t *testing.T) {
		mockClient := NewMockStoreClient()
		jobId := "test-job-multi-step"

		executor := kk_scheduler.NewTaskExecutor(
			kk_scheduler.WithSchedulerClient(mockClient),
			kk_scheduler.WithJobId(jobId),
		)

		step1Executed := false
		step2Executed := false

		executor.AddStep("step1", func(ctl *kk_scheduler.StepCtl) error {
			step1Executed = true
			return errors.New("step1 error")
		})

		executor.AddStep("step2", func(ctl *kk_scheduler.StepCtl) error {
			step2Executed = true
			ctl.Log("step2 done")
			return nil
		})

		err := executor.Run(context.Background())
		assert.NoError(t, err)
		assert.True(t, step1Executed)
		assert.True(t, step2Executed)

		logs := mockClient.GetStore().GetLogs(jobId)
		assert.Contains(t, logs, "❌step1 error")
		assert.Contains(t, logs, "step2 done")
		assert.Contains(t, logs, "[❌😭🤬Task Failed]")
	})
}
