package mock

import (
	"github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) AddTask(userID int64, task string) error {
	args := m.Called(userID, task)
	return args.Error(0)
}

func (m *MockDB) ListTasks(userID int64) ([]domain.Task, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockDB) CompleteTask(userID int64, taskID int64) error {
	args := m.Called(userID, taskID)
	return args.Error(0)
}
