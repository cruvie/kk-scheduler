package schedule

import (
	"github.com/cruvie/kk-scheduler/server/kk_scheduler"
)

type StoreDriver interface {
	JobList(serviceName string) ([]*kk_scheduler.PBJob, error)
	JobGet(serviceName, funcName string) (*kk_scheduler.PBJob, error)
	JobDelete(serviceName, funcName string) error
	JobPut(entry *kk_scheduler.PBJob) error

	ServiceList() ([]*kk_scheduler.PBRegisterService, error)
	ServicePut(service *kk_scheduler.PBRegisterService) error
	ServiceGet(serviceName string) (*kk_scheduler.PBRegisterService, error)
	ServiceDelete(serviceName string) error
}
