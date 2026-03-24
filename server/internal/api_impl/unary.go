package api_impl

import (
	"context"

	"gitee.com/cruvie/kk_go_kit/kk_grpc"
	"github.com/cruvie/kk-schedule/server/internal/api_handlers/job"
	"github.com/cruvie/kk-schedule/server/internal/api_handlers/service"
	"github.com/cruvie/kk-schedule/server/kk_schedule"
)

func (x *server) JobList(ctx context.Context, input *kk_schedule.JobList_Input) (*kk_schedule.JobList_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobList,
	)
}

func (x *server) JobGet(ctx context.Context, input *kk_schedule.JobGet_Input) (*kk_schedule.JobGet_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobGet,
	)
}

func (x *server) JobSetSpec(ctx context.Context, input *kk_schedule.JobSetSpec_Input) (*kk_schedule.JobSetSpec_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobSetSpec,
	)
}

func (x *server) JobEnable(ctx context.Context, input *kk_schedule.JobEnable_Input) (*kk_schedule.JobEnable_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobEnable,
	)
}

func (x *server) JobDisable(ctx context.Context, input *kk_schedule.JobDisable_Input) (*kk_schedule.JobDisable_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobDisable,
	)
}

func (x *server) JobPut(ctx context.Context, input *kk_schedule.JobPut_Input) (*kk_schedule.JobPut_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobPut,
	)
}

func (x *server) JobDelete(ctx context.Context, input *kk_schedule.JobDelete_Input) (*kk_schedule.JobDelete_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobDelete,
	)
}

func (x *server) JobTrigger(ctx context.Context, input *kk_schedule.JobTrigger_Input) (*kk_schedule.JobTrigger_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		job.NewApiJobTrigger,
	)
}

func (x *server) ServiceList(ctx context.Context, input *kk_schedule.ServiceList_Input) (*kk_schedule.ServiceList_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		service.NewApiServiceList,
	)
}

func (x *server) ServicePut(ctx context.Context, input *kk_schedule.ServicePut_Input) (*kk_schedule.ServicePut_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		service.NewApiServicePut,
	)
}

func (x *server) ServiceGet(ctx context.Context, input *kk_schedule.ServiceGet_Input) (*kk_schedule.ServiceGet_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		service.NewApiServiceGet,
	)
}

func (x *server) ServiceDelete(ctx context.Context, input *kk_schedule.ServiceDelete_Input) (*kk_schedule.ServiceDelete_Output, error) {
	return kk_grpc.GrpcHandler(
		ctx,
		input,
		service.NewApiServiceDelete,
	)
}
