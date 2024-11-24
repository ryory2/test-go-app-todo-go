// handler/task_test.go
package handler

import (
	"bytes"
	"encoding/json"
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

// TestCreateTask は CreateTask ハンドラーのテストです。
// このテストでは、モックリポジトリを使用してタスクの作成処理を検証します。
func TestCreateTask(t *testing.T) {
	// Gin をテストモードに設定
	gin.SetMode(gin.TestMode)

	// モックリポジトリのインスタンスを作成
	mockRepo := new(repository.MockTaskRepository)

	// バリデータのインスタンスを作成
	validate := validator.New()

	// ハンドラーのインスタンスを作成（モックリポジトリを注入）
	handler := NewTaskHandler(mockRepo, validate)

	// テスト用の Gin ルーターを作成し、ハンドラーを登録
	router := gin.Default()
	router.POST("/tasks", handler.CreateTask)

	// テストデータ（新しいタスク）の準備
	newTask := model.Task{
		Title:       "テストタスク",
		Description: "これはテスト用のタスクです。",
		DueDate:     time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	}

	// リクエストボディを JSON にエンコード
	jsonData, err := json.Marshal(newTask)
	assert.NoError(t, err) // エンコードにエラーがないことを確認

	// モックリポジトリの期待動作を設定
	mockRepo.On("CreateTask", mock.AnythingOfType("*model.Task")).Return(nil).Run(func(args mock.Arguments) {
		// CreateTask が呼ばれた際に、タスクの ID とタイムスタンプを設定
		task := args.Get(0).(*model.Task)
		task.ID = 1
		task.IsCompleted = false
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()
	})

	// テストリクエストを作成（POST /tasks）
	req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)                             // リクエスト作成にエラーがないことを確認
	req.Header.Set("Content-Type", "application/json") // ヘッダーを設定

	// レスポンスを記録するためのレスポンスライターを作成
	w := httptest.NewRecorder()

	// テストリクエストをルーターに送信
	router.ServeHTTP(w, req)

	// レスポンスのステータスコードが 201 Created であることを確認
	assert.Equal(t, http.StatusCreated, w.Code)

	// レスポンスボディを解析
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err) // JSON のパースにエラーがないことを確認

	// "data" フィールドが存在することを確認
	data, exists := response["data"].(map[string]interface{})
	assert.True(t, exists)

	// レスポンスデータの各フィールドを検証
	assert.Equal(t, float64(1), data["id"]) // JSONでは数値はfloat64として扱われる
	assert.Equal(t, newTask.Title, data["title"])
	assert.Equal(t, newTask.Description, data["description"])
	assert.Equal(t, newTask.DueDate.Format(time.RFC3339), data["due_date"])
	assert.Equal(t, false, data["is_completed"])

	// モックリポジトリが期待通りに呼び出されたことを確認
	mockRepo.AssertExpectations(t)
}
