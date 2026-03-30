package query

import (
	"context"
	"os"
	"testing"

	"gitee.com/cruvie/kk_go_kit/kk_env"
	"gitee.com/cruvie/kk_go_kit/kk_pg"
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/internal/models"
)

func init() {
	kk_env.SetEnv(kk_env.Env(os.Getenv("KK_Schedule")))
}

func TestGen(t *testing.T) {
	kk_env.SetEnv(kk_env.Env(os.Getenv("KK_Schedule")))
	stage := kk_stage.NewStage(context.Background(), "test")
	kk_pg.GenQuery(stage, &kk_pg.ConfigPG{DSN: kk_pg.PostgresDSN{
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "testpg",
		DBName:   "kk_scheduler",
		Schema:   "",
		SSLMode:  "disable",
		TimeZone: "UTC",
		Addition: nil,
	}},
		models.TaskExecution{},
		models.Job{},
		models.Service{},
	)
}
