package job

import (
	"gitee.com/cruvie/kk_go_kit/kk_grpc"
	"github.com/cruvie/kk-schedule/server/kk_schedule"
)

type ApiJobDelete struct {
	*kk_grpc.DefaultApi[kk_schedule.JobDelete_Input]
}

func NewApiJobDelete() *ApiJobDelete {
	return &ApiJobDelete{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.JobDelete_Input](),
	}
}

type ApiJobDisable struct {
	*kk_grpc.DefaultApi[kk_schedule.JobDisable_Input]
}

func NewApiJobDisable() *ApiJobDisable {
	return &ApiJobDisable{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.JobDisable_Input](),
	}
}

type ApiJobEnable struct {
	*kk_grpc.DefaultApi[kk_schedule.JobEnable_Input]
}

func NewApiJobEnable() *ApiJobEnable {
	return &ApiJobEnable{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.JobEnable_Input](),
	}
}

type ApiJobGet struct {
	*kk_grpc.DefaultApi[kk_schedule.JobGet_Input]
}

func NewApiJobGet() *ApiJobGet {
	return &ApiJobGet{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.JobGet_Input](),
	}
}

type ApiJobList struct {
	*kk_grpc.DefaultApi[kk_schedule.JobList_Input]
}

func NewApiJobList() *ApiJobList {
	return &ApiJobList{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.JobList_Input](),
	}
}

type ApiJobSetSpec struct {
	*kk_grpc.DefaultApi[kk_schedule.JobSetSpec_Input]
}

func NewApiJobSetSpec() *ApiJobSetSpec {
	return &ApiJobSetSpec{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.JobSetSpec_Input](),
	}
}

type ApiJobTrigger struct {
	*kk_grpc.DefaultApi[kk_schedule.JobTrigger_Input]
}

func NewApiJobTrigger() *ApiJobTrigger {
	return &ApiJobTrigger{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.JobTrigger_Input](),
	}
}

type ApiJobPut struct {
	*kk_grpc.DefaultApi[kk_schedule.JobPut_Input]
}

func NewApiJobPut() *ApiJobPut {
	return &ApiJobPut{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.JobPut_Input](),
	}
}
