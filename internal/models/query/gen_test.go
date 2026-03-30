package query

import (
	"context"
	"os"
	"testing"

	"gitee.com/cruvie/kk_go_kit/kk_env"
	"gitee.com/cruvie/kk_go_kit/kk_pg"
	"gitee.com/cruvie/kk_go_kit/kk_stage"
	"github.com/cruvie/kk-scheduler/internal/g_config"
	"github.com/cruvie/kk-scheduler/internal/models"
)

func init() {
	kk_env.SetEnv(kk_env.Env(os.Getenv("SS_Env")))
}

func TestGen(t *testing.T) {
	g_config.InitConfig()
	kk_env.SetEnv(kk_env.Env(os.Getenv("SS_Env")))
	stage := kk_stage.NewStage(context.Background(), "test")
	kk_pg.GenQuery(stage, g_config.Config.StorePG,
		models.TaskExecution{},
		models.Job{},
		models.Service{},
	)
}
