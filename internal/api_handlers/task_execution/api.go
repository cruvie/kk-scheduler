package task_execution

import (
	"gitee.com/cruvie/kk_go_kit/kk_grpc"
	"github.com/cruvie/kk-scheduler/kk_scheduler"
)

type ApiTaskCreate struct {
	*kk_grpc.DefaultApi[kk_scheduler.TaskCreate_Input]
}

func NewApiTaskCreate() *ApiTaskCreate {
	return &ApiTaskCreate{
		DefaultApi: kk_grpc.NewDefaultApi[kk_scheduler.TaskCreate_Input](),
	}
}

type ApiTaskUpdateStatus struct {
	*kk_grpc.DefaultApi[kk_scheduler.TaskUpdateStatus_Input]
}

func NewApiTaskUpdateStatus() *ApiTaskUpdateStatus {
	return &ApiTaskUpdateStatus{
		DefaultApi: kk_grpc.NewDefaultApi[kk_scheduler.TaskUpdateStatus_Input](),
	}
}

type ApiTaskAppendLog struct {
	*kk_grpc.DefaultApi[kk_scheduler.TaskAppendLog_Input]
}

func NewApiTaskAppendLog() *ApiTaskAppendLog {
	return &ApiTaskAppendLog{
		DefaultApi: kk_grpc.NewDefaultApi[kk_scheduler.TaskAppendLog_Input](),
	}
}
