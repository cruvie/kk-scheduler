package job

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/internal/schedule"
	"github.com/cruvie/kk-scheduler/kk_scheduler"
)

func (x *ApiJobDelete) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	err := schedule.GClient.JobDelete(x.In.GetId())
	return err
}

func (x *ApiJobDisable) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	err := schedule.GClient.JobDisable(x.In.GetId())
	return err
}

func (x *ApiJobEnable) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return schedule.GClient.JobEnable(x.In.GetId())
}

func (x *ApiJobGet) Service(stage *kk_stage.Stage) (*kk_scheduler.PBJob, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	job, err := schedule.GClient.JobGet(x.In.GetId())
	return job, err
}

func (x *ApiJobList) Service(stage *kk_stage.Stage) ([]*kk_scheduler.PBJob, error) {
	span := stage.StartTrace("Service")
	defer span.End()

	jobList, err := schedule.GClient.JobList(x.In.GetServiceName())
	return jobList, err
}

func (x *ApiJobSetSpec) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return schedule.GClient.JobSetSpec(x.In.GetId(), x.In.GetSpec())
}

func (x *ApiJobTrigger) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	return schedule.GClient.JobTrigger(x.In.GetId())
}

func (x *ApiJobPut) Service(stage *kk_stage.Stage) error {
	span := stage.StartTrace("Service")
	defer span.End()

	err := schedule.GClient.JobPut(x.In.GetJob())
	if err != nil {
		return err
	}

	return nil
}
