package store_driver

import (
	"github.com/cruvie/kk-scheduler/internal/models"
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
	TaskCreate(taskName string, status models.TaskExecutionStatus) error
	// TaskUpdateStatus updates the status
	TaskUpdateStatus(taskName string, status models.TaskExecutionStatus) error
	// TaskAppendLog append log to the execution record
	TaskAppendLog(taskName string, log string) error
}
