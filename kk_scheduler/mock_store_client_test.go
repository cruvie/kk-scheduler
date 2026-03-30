package kk_scheduler_test

import (
	"context"
	"sync"

	"google.golang.org/grpc"

	"github.com/cruvie/kk-scheduler/internal/store_driver"
	kk_scheduler "github.com/cruvie/kk-scheduler/kk_scheduler"
)

// MockStoreClient adapts store_driver.MockStore to KKScheduleClient interface
type MockStoreClient struct {
	store *store_driver.MockStore
	mu    sync.Mutex
	// idToTaskName maps task executor id to task name for status/log operations
	idToTaskName map[string]string
	// lastCreatedTaskName is the taskName of the most recently created task
	lastCreatedTaskName string
}

func NewMockStoreClient() *MockStoreClient {
	return &MockStoreClient{
		store:        store_driver.NewMockStore(),
		idToTaskName: make(map[string]string),
	}
}

func (m *MockStoreClient) GetStore() *store_driver.MockStore {
	return m.store
}

func (m *MockStoreClient) JobList(ctx context.Context, in *kk_scheduler.JobList_Input, opts ...grpc.CallOption) (*kk_scheduler.JobList_Output, error) {
	jobs, err := m.store.JobList(in.GetServiceName())
	if err != nil {
		return nil, err
	}
	out := &kk_scheduler.JobList_Output{}
	out.SetJobList(jobs)
	return out, nil
}

func (m *MockStoreClient) JobGet(ctx context.Context, in *kk_scheduler.JobGet_Input, opts ...grpc.CallOption) (*kk_scheduler.JobGet_Output, error) {
	job, err := m.store.JobGet(in.GetServiceName(), in.GetFuncName())
	if err != nil {
		return nil, err
	}
	out := &kk_scheduler.JobGet_Output{}
	out.SetJob(job)
	return out, nil
}

func (m *MockStoreClient) JobSetSpec(ctx context.Context, in *kk_scheduler.JobSetSpec_Input, opts ...grpc.CallOption) (*kk_scheduler.JobSetSpec_Output, error) {
	job, err := m.store.JobGet(in.GetServiceName(), in.GetFuncName())
	if err != nil {
		return nil, err
	}
	if job != nil {
		job.SetSpec(in.GetSpec())
	}
	return &kk_scheduler.JobSetSpec_Output{}, nil
}

func (m *MockStoreClient) JobEnable(ctx context.Context, in *kk_scheduler.JobEnable_Input, opts ...grpc.CallOption) (*kk_scheduler.JobEnable_Output, error) {
	job, err := m.store.JobGet(in.GetServiceName(), in.GetFuncName())
	if err != nil {
		return nil, err
	}
	if job != nil {
		job.SetEnabled(true)
	}
	return &kk_scheduler.JobEnable_Output{}, nil
}

func (m *MockStoreClient) JobDisable(ctx context.Context, in *kk_scheduler.JobDisable_Input, opts ...grpc.CallOption) (*kk_scheduler.JobDisable_Output, error) {
	job, err := m.store.JobGet(in.GetServiceName(), in.GetFuncName())
	if err != nil {
		return nil, err
	}
	if job != nil {
		job.SetEnabled(false)
	}
	return &kk_scheduler.JobDisable_Output{}, nil
}

func (m *MockStoreClient) JobPut(ctx context.Context, in *kk_scheduler.JobPut_Input, opts ...grpc.CallOption) (*kk_scheduler.JobPut_Output, error) {
	// Convert PBRegisterJob to PBJob
	registerJob := in.GetJob()
	pbJob := &kk_scheduler.PBJob{}
	pbJob.SetServiceName(registerJob.GetServiceName())
	pbJob.SetFuncName(registerJob.GetFuncName())
	pbJob.SetDescription(registerJob.GetDescription())
	err := m.store.JobPut(pbJob)
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.JobPut_Output{}, nil
}

func (m *MockStoreClient) JobDelete(ctx context.Context, in *kk_scheduler.JobDelete_Input, opts ...grpc.CallOption) (*kk_scheduler.JobDelete_Output, error) {
	err := m.store.JobDelete(in.GetServiceName(), in.GetFuncName())
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.JobDelete_Output{}, nil
}

func (m *MockStoreClient) JobTrigger(ctx context.Context, in *kk_scheduler.JobTrigger_Input, opts ...grpc.CallOption) (*kk_scheduler.JobTrigger_Output, error) {
	return &kk_scheduler.JobTrigger_Output{}, nil
}

func (m *MockStoreClient) ServiceList(ctx context.Context, in *kk_scheduler.ServiceList_Input, opts ...grpc.CallOption) (*kk_scheduler.ServiceList_Output, error) {
	services, err := m.store.ServiceList()
	if err != nil {
		return nil, err
	}
	out := &kk_scheduler.ServiceList_Output{}
	out.SetServiceList(services)
	return out, nil
}

func (m *MockStoreClient) ServicePut(ctx context.Context, in *kk_scheduler.ServicePut_Input, opts ...grpc.CallOption) (*kk_scheduler.ServicePut_Output, error) {
	err := m.store.ServicePut(in.GetService())
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.ServicePut_Output{}, nil
}

func (m *MockStoreClient) ServiceGet(ctx context.Context, in *kk_scheduler.ServiceGet_Input, opts ...grpc.CallOption) (*kk_scheduler.ServiceGet_Output, error) {
	service, err := m.store.ServiceGet(in.GetServiceName())
	if err != nil {
		return nil, err
	}
	out := &kk_scheduler.ServiceGet_Output{}
	out.SetService(service)
	return out, nil
}

func (m *MockStoreClient) ServiceDelete(ctx context.Context, in *kk_scheduler.ServiceDelete_Input, opts ...grpc.CallOption) (*kk_scheduler.ServiceDelete_Output, error) {
	err := m.store.ServiceDelete(in.GetServiceName())
	if err != nil {
		return nil, err
	}
	return &kk_scheduler.ServiceDelete_Output{}, nil
}

func (m *MockStoreClient) TaskCreate(ctx context.Context, in *kk_scheduler.TaskCreate_Input, opts ...grpc.CallOption) (*kk_scheduler.TaskCreate_Output, error) {
	taskName := in.GetJobId()
	m.store.TaskCreate(taskName)
	m.mu.Lock()
	m.idToTaskName[taskName] = taskName
	m.lastCreatedTaskName = taskName
	m.mu.Unlock()
	return &kk_scheduler.TaskCreate_Output{}, nil
}

func (m *MockStoreClient) TaskUpdateStatus(ctx context.Context, in *kk_scheduler.TaskUpdateStatus_Input, opts ...grpc.CallOption) (*kk_scheduler.TaskUpdateStatus_Output, error) {
	id := in.GetId()
	m.mu.Lock()
	taskName, ok := m.idToTaskName[id]
	if !ok && m.lastCreatedTaskName != "" {
		// First call with this id, map it to the last created task
		m.idToTaskName[id] = m.lastCreatedTaskName
		taskName = m.lastCreatedTaskName
	}
	m.mu.Unlock()
	if taskName != "" {
		m.store.TaskUpdateStatus(taskName, in.GetStatus())
	}
	return &kk_scheduler.TaskUpdateStatus_Output{}, nil
}

func (m *MockStoreClient) TaskAppendLog(ctx context.Context, in *kk_scheduler.TaskAppendLog_Input, opts ...grpc.CallOption) (*kk_scheduler.TaskAppendLog_Output, error) {
	id := in.GetId()
	m.mu.Lock()
	taskName, ok := m.idToTaskName[id]
	if !ok && m.lastCreatedTaskName != "" {
		// First call with this id, map it to the last created task
		m.idToTaskName[id] = m.lastCreatedTaskName
		taskName = m.lastCreatedTaskName
	}
	m.mu.Unlock()
	if taskName != "" {
		m.store.TaskAppendLog(taskName, in.GetLog())
	}
	return &kk_scheduler.TaskAppendLog_Output{}, nil
}
