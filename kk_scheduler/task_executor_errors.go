package kk_scheduler

import (
	"errors"
	"fmt"
)

// ErrStepPanic indicates a step panicked
var ErrStepPanic = errors.New("step panic")

// PanicError wraps panic information
type PanicError struct {
	PanicValue interface{}
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("step panicked: %v", e.PanicValue)
}
