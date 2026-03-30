package store_driver

import (
	"log/slog"
	"sync"
	"time"

	"github.com/cruvie/kk-scheduler/internal/models"
	"github.com/cruvie/kk-scheduler/internal/models/query"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresStore implements StoreDriver using PostgreSQL
type PostgresStore struct {
	db *gorm.DB
	q  *query.Query
	mu sync.Mutex
	// buffer for logs when connection fails
	logBuffer map[string][]string
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore(db *gorm.DB) *PostgresStore {
	q := query.Use(db)
	return &PostgresStore{
		db:        db,
		q:         q,
		logBuffer: make(map[string][]string),
	}
}

// Create creates a new task execution record
func (s *PostgresStore) Create(taskName string, status models.TaskExecutionStatus, startedAt string) (id string, err error) {
	id = uuid.New().String()
	execution := &models.TaskExecution{
		Id:        id,
		TaskName:  taskName,
		Status:    status,
		StartedAt: parseTime(startedAt),
		Log:       "",
	}

	if err = s.q.TaskExecution.Create(execution); err != nil {
		slog.Error("failed to create task execution", "err", err)
		return "", ErrStoreConnection
	}
	return id, nil
}

// UpdateStatus updates status and finished_at
func (s *PostgresStore) UpdateStatus(id string, status models.TaskExecutionStatus, finishedAt string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.q.TaskExecution.
		Where(s.q.TaskExecution.Id.Eq(id)).
		Update(s.q.TaskExecution.Status, status)
	if err != nil {
		return err
	}

	_, err = s.q.TaskExecution.
		Where(s.q.TaskExecution.Id.Eq(id)).
		Update(s.q.TaskExecution.FinishedAt, parseTime(finishedAt))
	return err
}

// AppendLog appends log message to the execution record
func (s *PostgresStore) AppendLog(id string, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Try direct append
	err := s.db.Exec(
		"UPDATE task_executions SET log = log || ? WHERE id = ?",
		message, id,
	).Error

	if err != nil {
		// Buffer the log for retry
		s.logBuffer[id] = append(s.logBuffer[id], message)
		slog.Warn("log write failed, buffering", "id", id, "err", err)
		return nil // Don't fail the task for log issues
	}

	// Flush any buffered logs
	if buffered := s.logBuffer[id]; len(buffered) > 0 {
		for _, msg := range buffered {
			s.db.Exec("UPDATE task_executions SET log = log || ? WHERE id = ?", msg, id)
		}
		delete(s.logBuffer, id)
	}
	return nil
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339Nano, s)
	return t
}
