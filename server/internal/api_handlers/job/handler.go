package job

import (
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-schedule/server/kk_schedule"
)

func (x *ApiJobDelete) Handler(stage *kk_stage.Stage) (*kk_schedule.JobDelete_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_schedule.JobDelete_Output{}, nil
}

func (x *ApiJobDisable) Handler(stage *kk_stage.Stage) (*kk_schedule.JobDisable_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_schedule.JobDisable_Output{}, nil
}

func (x *ApiJobEnable) Handler(stage *kk_stage.Stage) (*kk_schedule.JobEnable_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_schedule.JobEnable_Output{}, nil
}

func (x *ApiJobGet) Handler(stage *kk_stage.Stage) (*kk_schedule.JobGet_Output, error) {
	job, err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	out := &kk_schedule.JobGet_Output{}
	out.SetJob(job)
	return out, nil
}

func (x *ApiJobList) Handler(stage *kk_stage.Stage) (*kk_schedule.JobList_Output, error) {
	jobList, err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	out := &kk_schedule.JobList_Output{}
	out.SetJobList(jobList)
	return out, nil
}

func (x *ApiJobSetSpec) Handler(stage *kk_stage.Stage) (*kk_schedule.JobSetSpec_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_schedule.JobSetSpec_Output{}, nil
}

func (x *ApiJobTrigger) Handler(stage *kk_stage.Stage) (*kk_schedule.JobTrigger_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_schedule.JobTrigger_Output{}, nil
}

func (x *ApiJobPut) Handler(stage *kk_stage.Stage) (*kk_schedule.JobPut_Output, error) {
	err := x.Service(stage)
	if err != nil {
		return nil, err
	}
	return &kk_schedule.JobPut_Output{}, nil
}
