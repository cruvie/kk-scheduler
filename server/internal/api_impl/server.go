package api_impl

import (
	"gitee.com/cruvie/kk_go_kit/kk_grpc"
	"github.com/cruvie/kk-schedule/server/kk_schedule"
	"google.golang.org/grpc"
)

type server struct {
	kk_schedule.UnimplementedKKScheduleServer
}

func RegisterServer(grpcServer *grpc.Server) {
	kk_schedule.RegisterKKScheduleServer(grpcServer, &server{})
}

func RegisterFileDesc() {
	kk_grpc.GFileDescHub.RegisterFileDesc(kk_schedule.File_kk_schedule_service_proto)
}
