package task_execution

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/internal/schedule"
)

func (x *ApiTaskCreate) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return schedule.GClient.TaskCreate(x.In.GetJobId())
}

func (x *ApiTaskUpdateStatus) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return schedule.GClient.TaskUpdateStatus(x.In.GetId(), x.In.GetStatus())
}

func (x *ApiTaskAppendLog) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return schedule.GClient.TaskAppendLog(x.In.GetId(), x.In.GetLog())
}
