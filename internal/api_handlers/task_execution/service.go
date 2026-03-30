package task_execution

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/internal/models"
	"github.com/cruvie/kk-scheduler/internal/schedule"
)

func (x *ApiTaskCreate) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	status := models.TaskExecutionStatus(x.In.GetStatus())
	return schedule.GClient.TaskCreate(x.In.GetTaskName(), status)
}

func (x *ApiTaskUpdateStatus) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	status := models.TaskExecutionStatus(x.In.GetStatus())
	return schedule.GClient.TaskUpdateStatus(x.In.GetTaskName(), status)
}

func (x *ApiTaskAppendLog) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return schedule.GClient.TaskAppendLog(x.In.GetTaskName(), x.In.GetLog())
}
