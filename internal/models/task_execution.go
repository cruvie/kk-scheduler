package models

import (
	"time"

	"github.com/cruvie/kk-scheduler/kk_scheduler"
)

// TaskExecution 任务执行记录
type TaskExecution struct {
	Id         string                           `gorm:"primaryKey;column:id;type:uuid"`
	JobId      string                           `gorm:"column:job_id;type:uuid;not null"`
	Status     kk_scheduler.TaskExecutionStatus `gorm:"column:status;default:0;type:smallint;not null"`
	StartedAt  time.Time                        `gorm:"column:started_at;type:timestamp;not null"`
	FinishedAt time.Time                        `gorm:"column:finished_at;type:timestamp;not null"`
	Log        string                           `gorm:"column:log;type:text;not null"`
}

func (TaskExecution) TableName() string {
	return "task_executions"
}
