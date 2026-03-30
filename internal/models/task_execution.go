package models

import (
	"time"
)

// TaskExecution 任务执行记录
type TaskExecution struct {
	Id         string              `gorm:"primaryKey;column:id;type:uuid"`
	TaskName   string              `gorm:"column:task_name;type:text;not null"`
	Status     TaskExecutionStatus `gorm:"column:status;type:text;not null;index"`
	StartedAt  time.Time           `gorm:"column:started_at;type:timestamp;not null"`
	FinishedAt time.Time           `gorm:"column:finished_at;type:timestamp;not null"`
	Log        string              `gorm:"column:log;type:text;not null"`
}

func (TaskExecution) TableName() string {
	return "task_executions"
}

type TaskExecutionStatus string

// TaskExecutionStatus 任务执行状态常量
const (
	TaskExecutionStatusRunning   TaskExecutionStatus = "running"
	TaskExecutionStatusCompleted TaskExecutionStatus = "completed"
	TaskExecutionStatusFailed    TaskExecutionStatus = "failed"
	TaskExecutionStatusCancelled TaskExecutionStatus = "cancelled"
	TaskExecutionStatusTimeout   TaskExecutionStatus = "timeout"
)
