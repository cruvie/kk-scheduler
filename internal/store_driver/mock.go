package store_driver

import (
	"fmt"
	"sync"

	"github.com/cruvie/kk-scheduler/internal/models"
	"github.com/cruvie/kk-scheduler/kk_scheduler"
)

// MockStore is a thread-safe mock for testing
type MockStore struct {
	mu           sync.Mutex
	Records      map[string]*MockRecord
	Logs         map[string][]string
	CreatedCount int
	Jobs         map[string]*kk_scheduler.PBJob             // key: serviceName/funcName
	Services     map[string]*kk_scheduler.PBRegisterService // key: serviceName
}

type MockRecord struct {
	Id         string
	TaskName   string
	Status     models.TaskExecutionStatus
	StartedAt  string
	FinishedAt string
}

func NewMockStore() *MockStore {
	return &MockStore{
		Records:  make(map[string]*MockRecord),
		Logs:     make(map[string][]string),
		Jobs:     make(map[string]*kk_scheduler.PBJob),
		Services: make(map[string]*kk_scheduler.PBRegisterService),
	}
}

func (m *MockStore) TaskCreate(taskName string, status models.TaskExecutionStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CreatedCount++
	id := fmt.Sprintf("test-%d", m.CreatedCount)
	m.Records[id] = &MockRecord{
		Id:       id,
		TaskName: taskName,
		Status:   status,
	}
	m.Logs[id] = []string{}
	return nil
}

func (m *MockStore) TaskUpdateStatus(taskName string, status models.TaskExecutionStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range m.Records {
		if r.TaskName == taskName {
			r.Status = status
		}
	}
	return nil
}

func (m *MockStore) TaskAppendLog(taskName string, log string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, r := range m.Records {
		if r.TaskName == taskName {
			m.Logs[id] = append(m.Logs[id], log)
		}
	}
	return nil
}

func (m *MockStore) GetLogs(taskName string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, r := range m.Records {
		if r.TaskName == taskName {
			return m.Logs[id]
		}
	}
	return nil
}

func (m *MockStore) GetRecord(taskName string) *MockRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range m.Records {
		if r.TaskName == taskName {
			return r
		}
	}
	return nil
}

// Job methods
func (m *MockStore) jobKey(serviceName, funcName string) string {
	return fmt.Sprintf("%s/%s", serviceName, funcName)
}

func (m *MockStore) JobList(serviceName string) ([]*kk_scheduler.PBJob, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var jobs []*kk_scheduler.PBJob
	for _, job := range m.Jobs {
		if job.GetServiceName() == serviceName {
			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}

func (m *MockStore) JobGet(serviceName, funcName string) (*kk_scheduler.PBJob, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := m.jobKey(serviceName, funcName)
	return m.Jobs[key], nil
}

func (m *MockStore) JobDelete(serviceName, funcName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := m.jobKey(serviceName, funcName)
	delete(m.Jobs, key)
	return nil
}

func (m *MockStore) JobPut(entry *kk_scheduler.PBJob) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := m.jobKey(entry.GetServiceName(), entry.GetFuncName())
	m.Jobs[key] = entry
	return nil
}

// Service methods
func (m *MockStore) ServiceList() ([]*kk_scheduler.PBRegisterService, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var services []*kk_scheduler.PBRegisterService
	for _, svc := range m.Services {
		services = append(services, svc)
	}
	return services, nil
}

func (m *MockStore) ServiceGet(serviceName string) (*kk_scheduler.PBRegisterService, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Services[serviceName], nil
}

func (m *MockStore) ServicePut(service *kk_scheduler.PBRegisterService) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Services[service.GetServiceName()] = service
	return nil
}

func (m *MockStore) ServiceDelete(serviceName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Services, serviceName)
	return nil
}
