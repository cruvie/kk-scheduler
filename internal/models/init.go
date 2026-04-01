package models

import "gitee.com/cruvie/kk_go_kit/kk_pg"

func InitDB() {
	kk_pg.CreateTables(kk_pg.GormClient,
		TaskExecution{},
		Job{},
		Service{},
	)
}
