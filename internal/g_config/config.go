package g_config

import (
	"context"
	"log/slog"
	"os"

	"gitee.com/cruvie/kk_go_kit/kk_env"
	"gitee.com/cruvie/kk_go_kit/kk_pg"
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"gitee.com/cruvie/kk_go_kit/kk_time"
	"github.com/BurntSushi/toml"
	"github.com/cruvie/kk-scheduler/internal/models/query"
)

func init() {
	kk_env.SetEnv(kk_env.Env(os.Getenv("KK_Schedule")))
}

var Config config

type config struct {
	HttpPort   int
	GrpcPort   int
	WebPort    int
	StorePG    *kk_pg.ConfigPG
	configSlog *kk_stage.ConfigLog `toml:"-"`
}

func InitConfig() *kk_stage.Stage {
	data, err := os.ReadFile("/Users/cruvie/Documents/cruvie/kk-scheduler/config.toml")
	if err != nil {
		slog.Error("unable to read config.toml", "err", err)
		panic(err)
	}

	_, err = toml.Decode(string(data), &Config)
	if err != nil {
		slog.Error("unable to decode config.toml", "err", err)
		panic(err)
	}

	stage := kk_stage.NewStage(context.Background(), "kk-scheduler").SetStartTime(kk_time.NowUTCTime())
	{
		Config.configSlog = &kk_stage.ConfigLog{
			StartTime:  stage.StartTime,
			Lumberjack: kk_stage.DefaultLogConfig(kk_time.NowUTCTime(), "kk-scheduler"),
			Format:     kk_stage.FormatJSON,
		}
		Config.configSlog.Init()
	}
	{
		Config.StorePG.Init(stage)
		query.SetDefault(kk_pg.GormClient)
	}
	return stage
}

func CloseConfig() {
	Config.StorePG.Close()
	Config.configSlog.Close()
}
