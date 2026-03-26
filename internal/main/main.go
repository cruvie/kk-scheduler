package main

import (
	"context"
	"time"

	"gitee.com/cruvie/kk_go_kit/kk_server"
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"gitee.com/cruvie/kk_go_kit/kk_time"
	"github.com/cruvie/kk-scheduler/internal/g_config"
	"github.com/cruvie/kk-scheduler/internal/schedule"
)

var configSlog *kk_stage.ConfigLog

func main() {
	{
		g_config.InitConfig()
	}
	stage := kk_stage.NewStage(context.Background(), "kk-scheduler").SetStartTime(kk_time.NowUTCTime())
	{
		configSlog = &kk_stage.ConfigLog{
			StartTime:  stage.StartTime,
			Lumberjack: kk_stage.DefaultLogConfig(kk_time.NowUTCTime(), "kk-scheduler"),
			Format:     kk_stage.FormatJSON,
		}
		configSlog.Init()
		defer configSlog.Close()
	}
	kkServer := kk_server.NewKKServer(10*time.Second, stage)
	kkServer.Add("kk-scheduler", 0, schedule.NewScheduleServer())
	kkServer.Add("kk-scheduler-grpc", 0, NewGrpcServer(stage))
	kkServer.Add("kk-scheduler-http", 0, NewHttpServer(stage))
	kkServer.Add("kk-scheduler-web", 0, NewWebServer(stage))
	kkServer.ServeAndWait()
}
