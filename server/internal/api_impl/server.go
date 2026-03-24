package api_impl

import (
	"gitee.com/cruvie/kk_go_kit/kk_grpc"
	"github.com/cruvie/kk-scheduler/server/kk_scheduler"
	"google.golang.org/grpc"
)

type server struct {
	kk_scheduler.UnimplementedKKScheduleServer
}

func RegisterServer(grpcServer *grpc.Server) {
	kk_scheduler.RegisterKKScheduleServer(grpcServer, &server{})
}

func RegisterFileDesc() {
	kk_grpc.GFileDescHub.RegisterFileDesc(kk_scheduler.File_kk_scheduler_service_proto)
}
