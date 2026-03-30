package kk_scheduler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"gitee.com/cruvie/kk_go_kit/kk_id"
	"github.com/cruvie/kk-scheduler/kk_scheduler"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// getRealClient connects to the real gRPC scheduler server
func getRealClient(t *testing.T) (kk_scheduler.KKScheduleClient, func()) {
	conn, err := grpc.NewClient("127.0.0.1:8666",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	return kk_scheduler.NewKKScheduleClient(conn), func() { conn.Close() }
}

func TestTaskExecutor_RealServer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real server test in short mode")
	}

	client, cleanup := getRealClient(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	jobId := kk_id.GenUUID7()

	t.Run("executes steps and logs to server", func(t *testing.T) {
		executor := kk_scheduler.NewTaskExecutor(
			kk_scheduler.WithSchedulerClient(client),
			kk_scheduler.WithJobId(jobId),
		)

		stepExecuted := false
		executor.AddStep("real-step-1", func(ctl *kk_scheduler.StepCtl) error {
			stepExecuted = true
			ctl.Log("step executed on real server")
			return nil
		})

		executor.AddStep("real-step-2", func(ctl *kk_scheduler.StepCtl) error {
			ctl.Log("step 2 running")
			return nil
		})

		err := executor.Run(ctx)
		assert.NoError(t, err)
		assert.True(t, stepExecuted)
	})

	t.Run("handles step failure on real server", func(t *testing.T) {
		executor := kk_scheduler.NewTaskExecutor(
			kk_scheduler.WithSchedulerClient(client),
			kk_scheduler.WithJobId(kk_id.GenUUID7()),
		)

		executor.AddStep("failing-step", func(ctl *kk_scheduler.StepCtl) error {
			ctl.Log("about to fail")
			return errors.New("intentional failure for testing")
		})

		executor.AddStep("cleanup-step", func(ctl *kk_scheduler.StepCtl) error {
			ctl.Log("cleanup after failure")
			return nil
		})

		err := executor.Run(ctx)
		assert.NoError(t, err)
	})
}
