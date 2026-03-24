package service

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/server/internal/schedule"
	"github.com/cruvie/kk-scheduler/server/kk_scheduler"
)

func (x *ApiServiceDelete) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return schedule.GClient.ServiceDelete(x.In.GetServiceName())
}

func (x *ApiServiceGet) Service(stage *kk_stage.Stage) (*kk_scheduler.PBRegisterService, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	service, err := schedule.GClient.ServiceGet(x.In.GetServiceName())
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (x *ApiServiceList) Service(stage *kk_stage.Stage) ([]*kk_scheduler.PBRegisterService, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	service, err := schedule.GClient.ServiceList()
	return service, err
}

func (x *ApiServicePut) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return schedule.GClient.ServicePut(x.In.GetService())
}
