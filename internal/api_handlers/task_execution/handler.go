package task_execution

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/kk_scheduler"
)

func (x *ApiTaskCreate) Handler(stage *kk_stage.Stage) (*kk_scheduler.TaskCreate_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.TaskCreate_Output{}, nil
}

func (x *ApiTaskUpdateStatus) Handler(stage *kk_stage.Stage) (*kk_scheduler.TaskUpdateStatus_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.TaskUpdateStatus_Output{}, nil
}

func (x *ApiTaskAppendLog) Handler(stage *kk_stage.Stage) (*kk_scheduler.TaskAppendLog_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.TaskAppendLog_Output{}, nil
}
