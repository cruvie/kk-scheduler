package service

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-schedule/server/kk_schedule"
)

func (x *ApiServiceDelete) Handler(stage *kk_stage.Stage) (*kk_schedule.ServiceDelete_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_schedule.ServiceDelete_Output{}, nil
}

func (x *ApiServiceGet) Handler(stage *kk_stage.Stage) (*kk_schedule.ServiceGet_Output, error) {
	service, err := x.Service(stage)
	if err != nil {
		return nil, err
	}

	out := &kk_schedule.ServiceGet_Output{}
	out.SetService(service)
	return out, nil
}

func (x *ApiServiceList) Handler(stage *kk_stage.Stage) (*kk_schedule.ServiceList_Output, error) {
	service, err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	out := &kk_schedule.ServiceList_Output{}
	out.SetServiceList(service)
	return out, nil
}

func (x *ApiServicePut) Handler(stage *kk_stage.Stage) (*kk_schedule.ServicePut_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_schedule.ServicePut_Output{}, nil
}
