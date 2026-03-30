package kk_scheduler

// TaskExecutorOption configures TaskExecutor
type TaskExecutorOption func(*TaskExecutor)

// WithSchedulerClient sets the gRPC client for reporting
func WithSchedulerClient(client KKScheduleTaskExecutionClient) TaskExecutorOption {
	return func(t *TaskExecutor) {
		t.client = client
	}
}
