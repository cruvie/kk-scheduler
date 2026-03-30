package task_executor

import (
	"errors"
	"fmt"
)

var (
	// ErrTaskCancelled indicates task was manually cancelled
	ErrTaskCancelled = errors.New("task cancelled")
	// ErrTaskTimeout indicates task exceeded timeout
	ErrTaskTimeout = errors.New("task timeout")
	// ErrStepPanic indicates a step panicked
	ErrStepPanic = errors.New("step panic")
	// ErrStoreConnection indicates storage connection failure
	ErrStoreConnection = errors.New("store connection failure")
)

// PanicError wraps panic information
type PanicError struct {
	PanicValue interface{}
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("step panicked: %v", e.PanicValue)
}
