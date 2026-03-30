package store_driver

import (
	"github.com/cruvie/kk-scheduler/internal/g_config"
	"github.com/cruvie/kk-scheduler/kk_scheduler"
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

	// TaskCreate creates a new task execution record
	TaskCreate(jobId string) error
	// TaskUpdateStatus updates the status
	TaskUpdateStatus(id string, status kk_scheduler.TaskExecutionStatus) error
	// TaskAppendLog append log to the execution record
	TaskAppendLog(id string, log string) error
}

func NewStoreDriver() StoreDriver {
	switch g_config.Config.Store.Choose {
	case "PG":
		return NewPostgresStore()
	default:
		panic("store choose error")
	}
}
