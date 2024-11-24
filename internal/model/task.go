package model

import "time"

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" validate:"required,max=100"`
	Description string    `json:"description" validate:"omitempty,max=500"`
	DueDate     time.Time `json:"due_date" validate:"omitempty"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
