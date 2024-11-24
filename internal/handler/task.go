package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ryory2/test-go-app-todo-go/internal/model"
	"github.com/ryory2/test-go-app-todo-go/internal/repository"
)

// TaskHandler構造体
type TaskHandler struct {
	repo     repository.TaskRepository
	validate *validator.Validate
}

// NewTaskHandler関数
func NewTaskHandler(repo repository.TaskRepository, validate *validator.Validate) *TaskHandler {
	return &TaskHandler{
		repo:     repo,
		validate: validate,
	}
}

// GetTasksハンドラー
// HTTP: GET /tasks
func (h *TaskHandler) GetTasks(c *gin.Context) {
	// クエリパラメータの取得
	status := c.Query("status")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// クエリパラメータを整数に変換
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	// リポジトリを使用してタスクを取得
	tasks, total, err := h.repo.GetTasks(status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	// レスポンスを送信
	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"total": total,
	})
}

// CreateTaskハンドラー
// HTTP: POST /tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var input model.Task

	// リクエストボディをバインド
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// 入力値のバリデーション
	if err := h.validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// タスクを作成
	input.IsCompleted = false // 新規作成時は未完了とする
	if err := h.repo.CreateTask(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// 作成されたタスクを返す
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

// UpdateTaskハンドラー
// HTTP: PUT /tasks/{id}
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	// URLパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 既存のタスクを取得
	task, err := h.repo.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var input model.Task

	// リクエストボディをバインド
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// 入力値のバリデーション
	if err := h.validate.Struct(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// タスクのフィールドを更新
	task.Title = input.Title
	task.Description = input.Description
	task.DueDate = input.DueDate
	task.IsCompleted = input.IsCompleted
	task.UpdatedAt = time.Now()

	// タスクを更新
	if err := h.repo.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// 更新されたタスクを返す
	c.JSON(http.StatusOK, gin.H{"data": task})
}

// DeleteTaskハンドラー
// HTTP: DELETE /tasks/{id}
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// URLパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 既存のタスクを取得
	task, err := h.repo.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// タスクを削除
	if err := h.repo.DeleteTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	// 削除成功のレスポンスを送信
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// ToggleTaskハンドラー
// HTTP: PATCH /tasks/{id}/toggle
func (h *TaskHandler) ToggleTask(c *gin.Context) {
	// URLパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 既存のタスクを取得
	task, err := h.repo.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// タスクの完了状態をトグル
	if err := h.repo.ToggleTaskCompletion(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle task completion"})
		return
	}

	// 更新されたタスクを返す
	c.JSON(http.StatusOK, gin.H{"data": task})
}
