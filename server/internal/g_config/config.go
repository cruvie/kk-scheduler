package g_config

import (
	"log/slog"
	"os"

	"gitee.com/cruvie/kk_go_kit/kk_env"
	"gopkg.in/yaml.v3"
)

func init() {
	kk_env.SetEnv(kk_env.Env(os.Getenv("KK_Schedule")))
}

var Config config

type config struct {
	HttpPort  int `yaml:"HttpPort"`
	GrpcPort  int `yaml:"GrpcPort"`
	WebPort   int `yaml:"WebPort"`
	StoreEtcd struct {
		UserName  string   `yaml:"UserName"`
		Password  string   `yaml:"Password"`
		Endpoints []string `yaml:"Endpoints"`
	} `yaml:"StoreEtcd"`
}

func InitConfig() {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		slog.Error("unable to read config.yaml", "err", err)
		panic(err)
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		slog.Error("unable to unmarshal config.yaml", "err", err)
		panic(err)
	}
}
