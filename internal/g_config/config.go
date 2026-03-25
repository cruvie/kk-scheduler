package g_config

import (
	"log/slog"
	"os"

	"gitee.com/cruvie/kk_go_kit/kk_env"
	"github.com/BurntSushi/toml"
)

func init() {
	kk_env.SetEnv(kk_env.Env(os.Getenv("KK_Schedule")))
}

var Config config

type config struct {
	HttpPort  int
	GrpcPort  int
	WebPort   int
	StoreEtcd struct {
		UserName  string
		Password  string
		Endpoints []string
	}
}

func InitConfig() {
	data, err := os.ReadFile("config.toml")
	if err != nil {
		slog.Error("unable to read config.toml", "err", err)
		panic(err)
	}

	_, err = toml.Decode(string(data), &Config)
	if err != nil {
		slog.Error("unable to decode config.toml", "err", err)
		panic(err)
	}
}
