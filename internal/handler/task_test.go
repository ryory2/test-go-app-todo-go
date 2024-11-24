// internal/handler/task_test.go
package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ryory2/test-go-app-todo-go/internal/model"
	"github.com/ryory2/test-go-app-todo-go/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupTestHandler はテスト用の Gin エンジンとモックリポジトリをセットアップします。
func setupTestHandler(t *testing.T) (*gin.Engine, *repository.MockTaskRepository) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(repository.MockTaskRepository)
	validate := validator.New()
	handler := NewTaskHandler(mockRepo, validate)
	router := gin.Default()

	// エンドポイントの登録
	router.GET("/tasks", handler.GetTasks)
	router.POST("/tasks", handler.CreateTask)
	router.PUT("/tasks/:id", handler.UpdateTask)
	router.DELETE("/tasks/:id", handler.DeleteTask)
	router.PATCH("/tasks/:id/toggle", handler.ToggleTask)

	return router, mockRepo
}

// TestGetTasks は GetTasks ハンドラーの正常動作をテストします。
func TestGetTasks(t *testing.T) {
	router, mockRepo := setupTestHandler(t)

	// テストデータの準備
	tasks := []model.Task{
		{
			ID:          1,
			Title:       "タスク1",
			Description: "最初のタスク",
			IsCompleted: false,
		},
		{
			ID:          2,
			Title:       "タスク2",
			Description: "二つ目のタスク",
			IsCompleted: true,
		},
	}
	total := int64(len(tasks))

	// モックリポジトリの期待動作を設定
	mockRepo.On("GetTasks", "all", 10, 0).Return(tasks, total, nil)

	// テストリクエストを作成
	req, err := http.NewRequest(http.MethodGet, "/tasks?status=all&limit=10&offset=0", nil)
	assert.NoError(t, err)

	// レスポンスを記録するためのレスポンスライターを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	router.ServeHTTP(w, req)

	// レスポンスのステータスコードが 200 OK であることを確認
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスボディを解析
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// "data" フィールドが存在し、期待通りのタスク一覧が含まれていることを確認
	data, exists := response["data"].([]interface{})
	assert.True(t, exists)
	assert.Len(t, data, 2)

	// "total" フィールドが正しいことを確認
	assert.Equal(t, float64(total), response["total"])

	// モックリポジトリが期待通りに呼び出されたことを確認
	mockRepo.AssertExpectations(t)
}

// TestCreateTask は CreateTask ハンドラーの正常動作をテストします。
func TestCreateTask(t *testing.T) {
	router, mockRepo := setupTestHandler(t)

	// テストデータの準備
	newTask := model.Task{
		Title:       "テストタスク",
		Description: "これはテスト用のタスクです。",
		DueDate:     time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	}

	// リクエストボディの JSON エンコード
	jsonData, err := json.Marshal(newTask)
	assert.NoError(t, err)

	// モックリポジトリの期待動作を設定
	mockRepo.On("CreateTask", mock.AnythingOfType("*model.Task")).Return(nil).Run(func(args mock.Arguments) {
		task := args.Get(0).(*model.Task)
		task.ID = 1
		task.IsCompleted = false
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()
	})

	// テストリクエストを作成（POST /tasks）
	req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// レスポンスを記録するためのレスポンスライターを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	router.ServeHTTP(w, req)

	// レスポンスのステータスコードが 201 Created であることを確認
	assert.Equal(t, http.StatusCreated, w.Code)

	// レスポンスボディを解析
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// "data" フィールドが存在し、作成されたタスクが含まれていることを確認
	data, exists := response["data"].(map[string]interface{})
	assert.True(t, exists)
	assert.Equal(t, float64(1), data["id"]) // JSONでは数値はfloat64として扱われる
	assert.Equal(t, newTask.Title, data["title"])
	assert.Equal(t, newTask.Description, data["description"])
	assert.Equal(t, newTask.DueDate.Format(time.RFC3339), data["due_date"])
	assert.Equal(t, false, data["is_completed"])

	// モックリポジトリが期待通りに呼び出されたことを確認
	mockRepo.AssertExpectations(t)
}

// TestUpdateTask は UpdateTask ハンドラーの正常動作をテストします。
func TestUpdateTask(t *testing.T) {
	router, mockRepo := setupTestHandler(t)

	// 更新前のタスクデータ
	existingTask := &model.Task{
		ID:          1,
		Title:       "元のタスク",
		Description: "元の詳細",
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 更新後のタスクデータ（UpdatedAt は動的に設定されるため無視）
	updatedTask := &model.Task{
		ID:          1,
		Title:       "更新されたタスク",
		Description: "更新後の詳細",
		IsCompleted: true,
		CreatedAt:   existingTask.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	// モックリポジトリの期待動作を設定
	mockRepo.On("GetTaskByID", uint(1)).Return(existingTask, nil)
	// UpdatedTask の UpdatedAt は動的なので、特定の値ではなく、任意の *model.Task 型を受け取るように設定
	mockRepo.On("UpdateTask", mock.MatchedBy(func(t *model.Task) bool {
		return t.ID == updatedTask.ID &&
			t.Title == updatedTask.Title &&
			t.Description == updatedTask.Description &&
			t.IsCompleted == updatedTask.IsCompleted
	})).Return(nil)

	// テストリクエストの準備
	updateData := map[string]interface{}{
		"title":        "更新されたタスク",
		"description":  "更新後の詳細",
		"is_completed": true,
	}
	jsonData, err := json.Marshal(updateData)
	assert.NoError(t, err)

	// テストリクエストを作成（PUT /tasks/1）
	req, err := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// レスポンスを記録するためのレスポンスライターを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	router.ServeHTTP(w, req)

	// レスポンスのステータスコードが 200 OK であることを確認
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスボディを解析
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// "data" フィールドが存在し、更新されたタスクが含まれていることを確認
	data, exists := response["data"].(map[string]interface{})
	assert.True(t, exists)
	assert.Equal(t, float64(updatedTask.ID), data["id"])
	assert.Equal(t, updatedTask.Title, data["title"])
	assert.Equal(t, updatedTask.Description, data["description"])
	assert.Equal(t, updatedTask.IsCompleted, data["is_completed"])

	// モックリポジトリが期待通りに呼び出されたことを確認
	mockRepo.AssertExpectations(t)
}

// TestDeleteTask は DeleteTask ハンドラーの正常動作をテストします。
func TestDeleteTask(t *testing.T) {
	router, mockRepo := setupTestHandler(t)

	// 削除対象のタスクデータ
	existingTask := &model.Task{
		ID:          1,
		Title:       "削除するタスク",
		Description: "削除されるタスクの詳細",
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// モックリポジトリの期待動作を設定
	mockRepo.On("GetTaskByID", uint(1)).Return(existingTask, nil)
	mockRepo.On("DeleteTask", existingTask).Return(nil)

	// テストリクエストを作成（DELETE /tasks/1）
	req, err := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	assert.NoError(t, err)
	req.Header.Set("Accept", "application/json")

	// レスポンスを記録するためのレスポンスライターを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	router.ServeHTTP(w, req)

	// レスポンスのステータスコードが 200 OK であることを確認
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスボディを解析
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// "message" フィールドが存在し、削除成功メッセージが含まれていることを確認
	message, exists := response["message"].(string)
	assert.True(t, exists)
	assert.Equal(t, "Task deleted successfully", message)

	// モックリポジトリが期待通りに呼び出されたことを確認
	mockRepo.AssertExpectations(t)
}

// TestToggleTask は ToggleTask ハンドラーの正常動作をテストします。
func TestToggleTask(t *testing.T) {
	router, mockRepo := setupTestHandler(t)

	// 対象のタスクデータ（完了状態が false）
	existingTask := &model.Task{
		ID:          1,
		Title:       "トグルするタスク",
		Description: "トグルされるタスクの詳細",
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// トグル後のタスクデータ（完了状態が true）
	toggledTask := &model.Task{
		ID:          1,
		Title:       "トグルするタスク",
		Description: "トグルされるタスクの詳細", // Description は変更されない
		IsCompleted: true,
		CreatedAt:   existingTask.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	// モックリポジトリの期待動作を設定
	mockRepo.On("GetTaskByID", uint(1)).Return(existingTask, nil)
	mockRepo.On("ToggleTaskCompletion", existingTask).Return(nil).Run(func(args mock.Arguments) {
		// ToggleTaskCompletion が呼ばれた際に、IsCompleted フィールドを切り替える
		task := args.Get(0).(*model.Task)
		task.IsCompleted = !task.IsCompleted
		task.UpdatedAt = time.Now()
	})

	// テストリクエストを作成（PATCH /tasks/1/toggle）
	req, err := http.NewRequest(http.MethodPatch, "/tasks/1/toggle", nil)
	assert.NoError(t, err)
	req.Header.Set("Accept", "application/json")

	// レスポンスを記録するためのレスポンスライターを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	router.ServeHTTP(w, req)

	// レスポンスのステータスコードが 200 OK であることを確認
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスボディを解析
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// "data" フィールドが存在し、トグル後のタスクが含まれていることを確認
	data, exists := response["data"].(map[string]interface{})
	assert.True(t, exists)
	assert.Equal(t, float64(toggledTask.ID), data["id"])
	assert.Equal(t, toggledTask.Title, data["title"])
	assert.Equal(t, toggledTask.Description, data["description"]) // Description は変更されない
	assert.Equal(t, toggledTask.IsCompleted, data["is_completed"])

	// モックリポジトリが期待通りに呼び出されたことを確認
	mockRepo.AssertExpectations(t)
}

// TestGetTasks_RepositoryError は GetTasks ハンドラーのリポジトリエラー時のテストです。
func TestGetTasks_RepositoryError(t *testing.T) {
	router, mockRepo := setupTestHandler(t)

	// リポジトリがエラーを返すように設定（nil ではなく空のスライスを返す）
	mockRepo.On("GetTasks", "all", 10, 0).Return([]model.Task{}, int64(0), errors.New("database error"))

	// テストリクエストを作成
	req, err := http.NewRequest(http.MethodGet, "/tasks?status=all&limit=10&offset=0", nil)
	assert.NoError(t, err)

	// レスポンスを記録するためのレスポンスライターを作成
	w := httptest.NewRecorder()

	// リクエストをルーターに送信
	router.ServeHTTP(w, req)

	// レスポンスのステータスコードが 500 Internal Server Error であることを確認
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// レスポンスボディを解析
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// "error" フィールドが存在し、正しいエラーメッセージが含まれていることを確認
	errorMsg, exists := response["error"].(string)
	assert.True(t, exists)
	assert.Equal(t, "Failed to retrieve tasks", errorMsg)

	// モックリポジトリが期待通りに呼び出されたことを確認
	mockRepo.AssertExpectations(t)
}
