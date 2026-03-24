package schedule

import (
	"github.com/cruvie/kk-schedule/server/kk_schedule"
)

type StoreDriver interface {
	JobList(serviceName string) ([]*kk_schedule.PBJob, error)
	JobGet(serviceName, funcName string) (*kk_schedule.PBJob, error)
	JobDelete(serviceName, funcName string) error
	JobPut(entry *kk_schedule.PBJob) error

	ServiceList() ([]*kk_schedule.PBRegisterService, error)
	ServicePut(service *kk_schedule.PBRegisterService) error
	ServiceGet(serviceName string) (*kk_schedule.PBRegisterService, error)
	ServiceDelete(serviceName string) error
}
