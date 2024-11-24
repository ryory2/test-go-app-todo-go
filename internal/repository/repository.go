package repository

import (
	"github.com/ryory2/test-go-app-todo-go/internal/model"
	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(status string, limit, offset int) ([]model.Task, int64, error)
	CreateTask(task *model.Task) error
	GetTaskByID(id uint) (*model.Task, error)
	UpdateTask(task *model.Task) error
	DeleteTask(task *model.Task) error
	ToggleTaskCompletion(task *model.Task) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(status string, limit, offset int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64
	query := r.db.Model(&model.Task{})

	if status == "completed" {
		query = query.Where("is_completed = ?", true)
	} else if status == "pending" {
		query = query.Where("is_completed = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) CreateTask(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetTaskByID(id uint) (*model.Task, error) {
	var task model.Task
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) UpdateTask(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) DeleteTask(task *model.Task) error {
	return r.db.Delete(task).Error
}

func (r *taskRepository) ToggleTaskCompletion(task *model.Task) error {
	task.IsCompleted = !task.IsCompleted
	return r.db.Save(task).Error
}
