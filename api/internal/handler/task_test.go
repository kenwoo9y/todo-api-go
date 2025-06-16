package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kenwoo9y/todo-api-go/api/internal/entity"
	"github.com/kenwoo9y/todo-api-go/api/pkg/common"
)

// MockTaskRepository は repository.TaskRepository のモック実装
type MockTaskRepository struct {
	createFunc       func(ctx context.Context, task *entity.Task) error
	getAllFunc       func(ctx context.Context) ([]entity.Task, error)
	getByIDFunc      func(ctx context.Context, id int64) (*entity.Task, error)
	getByOwnerIDFunc func(ctx context.Context, ownerID int64) ([]entity.Task, error)
	updateFunc       func(ctx context.Context, task *entity.Task) error
	deleteFunc       func(ctx context.Context, id int64) error
}

func (m *MockTaskRepository) Create(ctx context.Context, task *entity.Task) error {
	return m.createFunc(ctx, task)
}

func (m *MockTaskRepository) GetAll(ctx context.Context) ([]entity.Task, error) {
	return m.getAllFunc(ctx)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id int64) (*entity.Task, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *MockTaskRepository) GetByOwnerID(ctx context.Context, ownerID int64) ([]entity.Task, error) {
	return m.getByOwnerIDFunc(ctx, ownerID)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *entity.Task) error {
	return m.updateFunc(ctx, task)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int64) error {
	return m.deleteFunc(ctx, id)
}

func TestTaskHandler_Create(t *testing.T) {
	now := time.Now().Format("2006-01-02")
	tests := []struct {
		name           string
		requestBody    CreateTaskRequest
		mockSetup      func(*MockTaskRepository)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "正常系：タスク作成成功",
			requestBody: CreateTaskRequest{
				Title:       "テストタスク",
				Description: "テストの説明",
				DueDate:     now,
				Status:      "Todo",
				OwnerID:     1,
			},
			mockSetup: func(m *MockTaskRepository) {
				m.createFunc = func(ctx context.Context, task *entity.Task) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "異常系：リポジトリエラー",
			requestBody: CreateTaskRequest{
				Title:       "テストタスク",
				Description: "テストの説明",
				DueDate:     now,
				Status:      "Todo",
				OwnerID:     1,
			},
			mockSetup: func(m *MockTaskRepository) {
				m.createFunc = func(ctx context.Context, task *entity.Task) error {
					return errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name: "異常系：不正な日付フォーマット",
			requestBody: CreateTaskRequest{
				Title:       "テストタスク",
				Description: "テストの説明",
				DueDate:     "2025/06/15", // 不正なフォーマット
				Status:      "Todo",
				OwnerID:     1,
			},
			mockSetup: func(m *MockTaskRepository) {
				m.createFunc = func(ctx context.Context, task *entity.Task) error {
					return nil
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTaskRepository{}
			tt.mockSetup(mockRepo)

			handler := NewTaskHandler(mockRepo)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Create(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response entity.Task
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response.Title != tt.requestBody.Title {
					t.Errorf("expected title %s, got %s", tt.requestBody.Title, response.Title)
				}
				if response.Description != tt.requestBody.Description {
					t.Errorf("expected description %s, got %s", tt.requestBody.Description, response.Description)
				}
				if response.Status != entity.TaskStatus(tt.requestBody.Status) {
					t.Errorf("expected status %s, got %s", tt.requestBody.Status, response.Status)
				}
				if response.OwnerID != tt.requestBody.OwnerID {
					t.Errorf("expected owner ID %d, got %d", tt.requestBody.OwnerID, response.OwnerID)
				}
			}
		})
	}
}

func TestTaskHandler_GetByID(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name           string
		taskID         int64
		mockSetup      func(*MockTaskRepository)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "正常系：タスク取得成功",
			taskID: 1,
			mockSetup: func(m *MockTaskRepository) {
				m.getByIDFunc = func(ctx context.Context, id int64) (*entity.Task, error) {
					return &entity.Task{
						ID:          1,
						Title:       "テストタスク",
						Description: "テストの説明",
						DueDate:     now.Format("2006-01-02"),
						Status:      entity.TaskStatusTodo,
						OwnerID:     1,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "異常系：タスクが見つからない",
			taskID: 999,
			mockSetup: func(m *MockTaskRepository) {
				m.getByIDFunc = func(ctx context.Context, id int64) (*entity.Task, error) {
					return nil, common.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTaskRepository{}
			tt.mockSetup(mockRepo)

			handler := NewTaskHandler(mockRepo)
			req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
			w := httptest.NewRecorder()

			handler.GetByID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response entity.Task
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response.ID != tt.taskID {
					t.Errorf("expected task ID %d, got %d", tt.taskID, response.ID)
				}
			}
		})
	}
}

func TestTaskHandler_Update(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name           string
		taskID         int64
		requestBody    UpdateTaskRequest
		mockSetup      func(*MockTaskRepository)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "正常系：タスク更新成功",
			taskID: 1,
			requestBody: UpdateTaskRequest{
				Title:       taskStringPtr("更新されたタスク"),
				Description: taskStringPtr("更新された説明"),
				Status:      taskStringPtr("Done"),
			},
			mockSetup: func(m *MockTaskRepository) {
				m.getByIDFunc = func(ctx context.Context, id int64) (*entity.Task, error) {
					return &entity.Task{
						ID:          1,
						Title:       "テストタスク",
						Description: "テストの説明",
						DueDate:     now.Format("2006-01-02"),
						Status:      entity.TaskStatusTodo,
						OwnerID:     1,
					}, nil
				}
				m.updateFunc = func(ctx context.Context, task *entity.Task) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "異常系：更新フィールドが空",
			taskID: 1,
			requestBody: UpdateTaskRequest{
				Title:       nil,
				Description: nil,
				Status:      nil,
			},
			mockSetup: func(m *MockTaskRepository) {
				m.getByIDFunc = func(ctx context.Context, id int64) (*entity.Task, error) {
					return &entity.Task{
						ID:          1,
						Title:       "テストタスク",
						Description: "テストの説明",
						DueDate:     now.Format("2006-01-02"),
						Status:      entity.TaskStatusTodo,
						OwnerID:     1,
					}, nil
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockTaskRepository{}
			tt.mockSetup(mockRepo)

			handler := NewTaskHandler(mockRepo)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/tasks/1", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Update(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response entity.Task
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if tt.requestBody.Title != nil && response.Title != *tt.requestBody.Title {
					t.Errorf("expected title %s, got %s", *tt.requestBody.Title, response.Title)
				}
				if tt.requestBody.Description != nil && response.Description != *tt.requestBody.Description {
					t.Errorf("expected description %s, got %s", *tt.requestBody.Description, response.Description)
				}
				if tt.requestBody.Status != nil && response.Status != entity.TaskStatus(*tt.requestBody.Status) {
					t.Errorf("expected status %s, got %s", *tt.requestBody.Status, response.Status)
				}
			}
		})
	}
}

// ヘルパー関数
func taskStringPtr(s string) *string {
	return &s
}
