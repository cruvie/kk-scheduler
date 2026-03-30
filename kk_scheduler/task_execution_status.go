package kk_scheduler

// TaskExecutionStatus 任务执行状态
type TaskExecutionStatus string

// TaskExecutionStatus 任务执行状态常量
const (
	TaskExecutionStatusRunning   TaskExecutionStatus = "running"
	TaskExecutionStatusCompleted TaskExecutionStatus = "completed"
	TaskExecutionStatusFailed    TaskExecutionStatus = "failed"
)
