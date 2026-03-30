package task_executor

import (
	"time"

	"github.com/cruvie/kk-scheduler/internal/store_driver"
	"github.com/cruvie/kk-scheduler/kk_scheduler"
	"gorm.io/gorm"
)

// TaskOption configures TaskExecutor
type TaskOption func(*TaskExecutor)

// WithTimeout sets the maximum execution time
func WithTimeout(d time.Duration) TaskOption {
	return func(t *TaskExecutor) {
		t.timeout = d
	}
}

// WithStore sets a custom store driver
func WithStore(store store_driver.StoreDriver) TaskOption {
	return func(t *TaskExecutor) {
		t.store = store
	}
}

// GRPCClientBuilder is a function that creates a gRPC client for termination check
// The builder function should return a fresh client each call
type GRPCClientBuilder func() (kk_scheduler.KKScheduleTriggerClient, error)

// WithGRPCClient sets the gRPC client builder for termination check
func WithGRPCClient(builder GRPCClientBuilder) TaskOption {
	return func(t *TaskExecutor) {
		t.grpcClientBuilder = builder
	}
}

// WithDB sets the database connection for default PostgresStore
func WithDB(db *gorm.DB) TaskOption {
	return func(t *TaskExecutor) {
		t.store = store_driver.NewPostgresStore(db)
	}
}
