package task_execution

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/internal/scheduler"
)

func (x *ApiTaskCreate) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.TaskCreate(x.In)
}

func (x *ApiTaskUpdateStatus) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.TaskUpdateStatus(x.In.GetId(), x.In.GetStatus())
}

func (x *ApiTaskAppendLog) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return scheduler.GClient.TaskAppendLog(x.In.GetId(), x.In.GetLog())
}
