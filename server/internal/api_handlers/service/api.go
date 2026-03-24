package service

import (
	"gitee.com/cruvie/kk_go_kit/kk_grpc"

	"github.com/cruvie/kk-schedule/server/kk_schedule"
)

type ApiServiceDelete struct {
	*kk_grpc.DefaultApi[kk_schedule.ServiceDelete_Input]
}

func NewApiServiceDelete() *ApiServiceDelete {
	return &ApiServiceDelete{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.ServiceDelete_Input](),
	}
}

type ApiServiceGet struct {
	*kk_grpc.DefaultApi[kk_schedule.ServiceGet_Input]
}

func NewApiServiceGet() *ApiServiceGet {
	return &ApiServiceGet{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.ServiceGet_Input](),
	}
}

type ApiServiceList struct {
	*kk_grpc.DefaultApi[kk_schedule.ServiceList_Input]
}

func NewApiServiceList() *ApiServiceList {
	return &ApiServiceList{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.ServiceList_Input](),
	}
}

type ApiServicePut struct {
	*kk_grpc.DefaultApi[kk_schedule.ServicePut_Input]
}

func NewApiServicePut() *ApiServicePut {
	return &ApiServicePut{
		DefaultApi: kk_grpc.NewDefaultApi[kk_schedule.ServicePut_Input](),
	}
}
