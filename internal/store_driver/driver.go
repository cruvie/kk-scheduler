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

	// Create creates a new task execution record and returns the ID
	Create(taskName string, status models.TaskExecutionStatus, startedAt string) (id string, err error)
	// UpdateStatus updates the status and finished_at time
	UpdateStatus(id string, status models.TaskExecutionStatus, finishedAt string) error
	// AppendLog appends log message to the execution record
	AppendLog(id string, message string) error
}
