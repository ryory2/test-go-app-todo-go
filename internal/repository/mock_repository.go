// internal/repository/mock_repository.go
package repository

import (
	"github.com/ryory2/test-go-app-todo-go/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository は TaskRepository インターフェースのモック実装です
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetTasks(status string, limit, offset int) ([]model.Task, int64, error) {
	args := m.Called(status, limit, offset)
	return args.Get(0).([]model.Task), args.Get(1).(int64), args.Error(2)
}

func (m *MockTaskRepository) CreateTask(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetTaskByID(id uint) (*model.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockTaskRepository) UpdateTask(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteTask(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) ToggleTaskCompletion(task *model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}
