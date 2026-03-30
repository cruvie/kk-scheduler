package store_driver

import (
	"fmt"
	"sync"

	"github.com/cruvie/kk-scheduler/internal/models"
)

// MockStore is a thread-safe mock for testing
type MockStore struct {
	mu           sync.Mutex
	Records      map[string]*MockRecord
	Logs         map[string][]string
	CreatedCount int
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
		Records: make(map[string]*MockRecord),
		Logs:    make(map[string][]string),
	}
}

func (m *MockStore) Create(taskName string, status models.TaskExecutionStatus, startedAt string) (id string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CreatedCount++
	id = fmt.Sprintf("test-%d", m.CreatedCount)
	m.Records[id] = &MockRecord{
		Id:        id,
		TaskName:  taskName,
		Status:    status,
		StartedAt: startedAt,
	}
	m.Logs[id] = []string{}
	return id, nil
}

func (m *MockStore) UpdateStatus(id string, status models.TaskExecutionStatus, finishedAt string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r, ok := m.Records[id]; ok {
		r.Status = status
		r.FinishedAt = finishedAt
	}
	return nil
}

func (m *MockStore) AppendLog(id string, message string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logs[id] = append(m.Logs[id], message)
	return nil
}

func (m *MockStore) GetLogs(id string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Logs[id]
}

func (m *MockStore) GetRecord(id string) *MockRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Records[id]
}
