package schedule

import (
	"errors"
	"slices"
	"strings"

	"github.com/cruvie/kk-scheduler/internal/models"
	"github.com/cruvie/kk-scheduler/internal/store_driver"
	"github.com/cruvie/kk-scheduler/kk_scheduler"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GClient *Client

type Client struct {
	cron   *cron.Cron
	storer store_driver.StoreDriver
}

func InitGClient(cfg *Config) {
	cfg.check()
	c := &Client{
		cron:   cron.New(cfg.Opts...),
		storer: cfg.StoreDriver,
	}
	GClient = c
	GClient.initJob()
}

func (x *Client) initJob() {
	jobList, err := x.JobList("")
	if err != nil {
		panic(err)
	}
	for _, job := range jobList {
		if !job.GetEnabled() {
			continue
		}
		err := x.JobEnable(job.GetServiceName(), job.GetFuncName())
		if err != nil {
			panic(err)
		}
	}
}

func (x *Client) Start() {
	x.cron.Start()
}

func (x *Client) Close() {
	x.cron.Stop()
}

func (x *Client) JobPut(jobs ...*kk_scheduler.PBRegisterJob) error {
	for _, job := range jobs {
		err := job.Check()
		if err != nil {
			panic(err)
		}
		{ // check service exist
			_, err = x.storer.ServiceGet(job.GetServiceName())
			if err != nil {
				return err
			}
		}

		newEntry := &kk_scheduler.PBJob{}
		newEntry.SetEntryID(0)
		newEntry.SetEnabled(false)
		newEntry.SetNext(nil)
		newEntry.SetPrev(nil)
		newEntry.SetSpec("")
		newEntry.SetServiceName(job.GetServiceName())
		newEntry.SetDescription(job.GetDescription())
		newEntry.SetFuncName(job.GetFuncName())

		entry, err := x.storer.JobGet(job.GetServiceName(), job.GetFuncName())
		if err != nil && !errors.Is(err, kk_scheduler.ErrJobNotFount) {
			return err
		}

		if entry != nil {
			newEntry.SetEntryID(entry.GetEntryID())
			newEntry.SetEnabled(entry.GetEnabled())
			newEntry.SetNext(entry.GetNext())
			newEntry.SetPrev(entry.GetPrev())
			newEntry.SetSpec(entry.GetSpec())
		}
		err = x.storer.JobPut(newEntry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *Client) JobList(serviceName string) ([]*kk_scheduler.PBJob, error) {
	entries := x.cron.Entries()
	entryList, err := x.storer.JobList(serviceName)
	if err != nil {
		return nil, err
	}
	var hasSpecEntryList []int32
	var pbJobs []*kk_scheduler.PBJob
	for _, entry := range entries {
		find, b := lo.Find(entryList, func(item *kk_scheduler.PBJob) bool {
			return item.GetEntryID() == int32(entry.ID)
		})
		if !b {
			continue
		} else {
			hasSpecEntryList = append(hasSpecEntryList, find.GetEntryID())
		}
		job := &kk_scheduler.PBJob{}
		job.SetEntryID(find.GetEntryID())
		job.SetEnabled(find.GetEnabled())
		job.SetNext(find.GetNext())
		job.SetPrev(find.GetPrev())
		job.SetSpec(find.GetSpec())
		job.SetDescription(find.GetDescription())
		job.SetFuncName(find.GetFuncName())
		job.SetServiceName(find.GetServiceName())
		pbJobs = append(pbJobs, job)
	}

	noSpecJobList := lo.Filter(entryList, func(item *kk_scheduler.PBJob, index int) bool {
		_, b := lo.Find(hasSpecEntryList, func(id int32) bool {
			return id == item.GetEntryID()
		})
		return !b
	})
	pbJobs = append(pbJobs, noSpecJobList...)
	// sort by serviceName
	slices.SortFunc(pbJobs, func(a, b *kk_scheduler.PBJob) int {
		return strings.Compare(a.GetServiceName(), b.GetServiceName())
	})
	return pbJobs, nil
}

func (x *Client) JobGet(serviceName, funcName string) (*kk_scheduler.PBJob, error) {
	entry, err := x.storer.JobGet(serviceName, funcName)
	if err != nil {
		return nil, err
	}
	cEntry := x.cron.Entry(cron.EntryID(entry.GetEntryID()))

	if cEntry.Valid() {
		entry.SetNext(timestamppb.New(cEntry.Next))
		entry.SetPrev(timestamppb.New(cEntry.Prev))
	}

	return entry, nil
}

func (x *Client) JobSetSpec(serviceName, funcName string, spec string) error {
	job, err := x.storer.JobGet(serviceName, funcName)
	if err != nil {
		return err
	}
	if job.GetSpec() == spec {
		return nil
	}
	job.SetSpec(spec)

	err = x.storer.JobPut(job)
	if job.GetEnabled() {
		err = x.JobEnable(serviceName, funcName)
		if err != nil {
			return err
		}
	}

	return err
}

func (x *Client) JobEnable(serviceName string, funcName string) error {
	entry, err := x.storer.JobGet(serviceName, funcName)
	if err != nil {
		return err
	}
	if entry.GetSpec() == "" {
		return kk_scheduler.ErrSpecIsEmpty
	}
	err = x.JobDisable(serviceName, funcName)
	if err != nil {
		return err
	}
	service, err := x.storer.ServiceGet(serviceName)
	if err != nil {
		return err
	}

	entryID, err := x.cron.AddFunc(entry.GetSpec(), triggerFunc(service, funcName))
	if err != nil {
		return err
	}

	entry.SetEntryID(int32(entryID))
	entry.SetEnabled(true)

	err = x.storer.JobPut(entry)
	return err
}

func (x *Client) JobDisable(serviceName, funcName string) error {
	entry, err := x.storer.JobGet(serviceName, funcName)
	if err != nil {
		return err
	}
	x.cron.Remove(cron.EntryID(entry.GetEntryID()))

	entry.SetEnabled(false)
	entry.SetEntryID(0)

	err = x.storer.JobPut(entry)
	return err
}

func (x *Client) JobDelete(serviceName, funcName string) error {
	// disable job
	err := x.JobDisable(serviceName, funcName)
	if err != nil {
		return err
	}
	return x.storer.JobDelete(serviceName, funcName)
}

// JobTrigger triggers a job manually
func (x *Client) JobTrigger(serviceName, funcName string) error {
	service, err := x.storer.ServiceGet(serviceName)
	if err != nil {
		return err
	}

	// Trigger the job function directly
	triggerFunc(service, funcName)()
	return nil
}

func (x *Client) ServiceList() ([]*kk_scheduler.PBRegisterService, error) {
	return x.storer.ServiceList()
}

func (x *Client) ServicePut(service *kk_scheduler.PBRegisterService) error {
	return x.storer.ServicePut(service)
}

func (x *Client) ServiceGet(serviceName string) (*kk_scheduler.PBRegisterService, error) {
	return x.storer.ServiceGet(serviceName)
}

func (x *Client) ServiceDelete(serviceName string) error {
	// check no job in service
	jobList, err := x.JobList(serviceName)
	if err != nil {
		return err
	}
	if len(jobList) > 0 {
		return kk_scheduler.ErrServiceHasJob
	}
	return x.storer.ServiceDelete(serviceName)
}

// TaskCreate creates a new task execution record
func (x *Client) TaskCreate(taskName string, status models.TaskExecutionStatus) error {
	return x.storer.TaskCreate(taskName, status)
}

// TaskUpdateStatus updates the task execution status
func (x *Client) TaskUpdateStatus(taskName string, status models.TaskExecutionStatus) error {
	return x.storer.TaskUpdateStatus(taskName, status)
}

// TaskAppendLog appends log to the task execution record
func (x *Client) TaskAppendLog(taskName string, log string) error {
	return x.storer.TaskAppendLog(taskName, log)
}
