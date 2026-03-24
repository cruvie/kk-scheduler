package internal

import (
	"testing"

	"gitee.com/cruvie/kk_go_kit/kk_grpc/grpc_api_gen"
	"gitee.com/cruvie/kk_go_kit/kk_system"
	"github.com/cruvie/kk-schedule/server/kk_schedule"
)

func TestName(t *testing.T) {
	kk_system.TerminatePort(3000)
}

func TestGeneratePermissionApi(t *testing.T) {
	//apiGroupModel := grpc_api_gen.ApiGroupModel{
	//	AdditionImports: `
	//`,
	//	TargetPath: "./api_handlers",
	//}

	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.JobList{},
	//})
	//
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.JobGet{},
	//})
	//
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.JobSetSpec{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.JobEnable{},
	//})
	//
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.JobDisable{},
	//})
	//
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.JobPut{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.JobDelete{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.JobTrigger{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.ServiceList{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.ServicePut{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.ServiceGet{},
	//})
	//grpc_api_gen.GenerateHandler(apiGroupModel, grpc_api_gen.ApiModel{
	//	ApiPtr: &kk_schedule.ServiceDelete{},
	//})
}

func TestGenImpl(t *testing.T) {
	grpc_api_gen.GenerateImpl(
		grpc_api_gen.GenerateImplInput{
			ServerName:      "KKSchedule",
			Methods:         kk_schedule.KKSchedule_ServiceDesc.Methods,
			ApiDefPkgPath:   "github.com/cruvie/kk-schedule/server/internal/api_def",
			HandlersPkgPath: "github.com/cruvie/kk-schedule/server/internal/api_handlers",
		})
}
