package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kenwoo9y/todo-api-go/api/internal/entity"
	"github.com/kenwoo9y/todo-api-go/api/pkg/common"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	createFunc        func(ctx context.Context, user *entity.User) error
	getAllFunc        func(ctx context.Context) ([]entity.User, error)
	getByIDFunc       func(ctx context.Context, id int64) (*entity.User, error)
	getByUsernameFunc func(ctx context.Context, username string) (*entity.User, error)
	updateFunc        func(ctx context.Context, user *entity.User) error
	deleteFunc        func(ctx context.Context, id int64) error
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	return m.createFunc(ctx, user)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	return m.getAllFunc(ctx)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return m.getByUsernameFunc(ctx, username)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	return m.updateFunc(ctx, user)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	return m.deleteFunc(ctx, id)
}

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateUserRequest
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Success: User creation succeeds",
			requestBody: CreateUserRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
			},
			mockSetup: func(m *MockUserRepository) {
				m.createFunc = func(ctx context.Context, user *entity.User) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "Error: Repository error",
			requestBody: CreateUserRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				FirstName: "Test",
				LastName:  "User",
			},
			mockSetup: func(m *MockUserRepository) {
				m.createFunc = func(ctx context.Context, user *entity.User) error {
					return errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.mockSetup(mockRepo)

			handler := NewUserHandler(mockRepo)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Create(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response entity.User
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response.Username != tt.requestBody.Username {
					t.Errorf("expected username %s, got %s", tt.requestBody.Username, response.Username)
				}
				if response.Email != tt.requestBody.Email {
					t.Errorf("expected email %s, got %s", tt.requestBody.Email, response.Email)
				}
			}
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "Success: User retrieval succeeds",
			userID: 1,
			mockSetup: func(m *MockUserRepository) {
				m.getByIDFunc = func(ctx context.Context, id int64) (*entity.User, error) {
					return &entity.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						FirstName: "Test",
						LastName:  "User",
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "Error: User not found",
			userID: 999,
			mockSetup: func(m *MockUserRepository) {
				m.getByIDFunc = func(ctx context.Context, id int64) (*entity.User, error) {
					return nil, common.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.mockSetup(mockRepo)

			handler := NewUserHandler(mockRepo)
			req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
			w := httptest.NewRecorder()

			handler.GetByID(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response entity.User
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response.ID != tt.userID {
					t.Errorf("expected user ID %d, got %d", tt.userID, response.ID)
				}
			}
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		requestBody    UpdateUserRequest
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "Success: User update succeeds",
			userID: 1,
			requestBody: UpdateUserRequest{
				Username: stringPtr("updateduser"),
				Email:    stringPtr("updated@example.com"),
			},
			mockSetup: func(m *MockUserRepository) {
				m.getByIDFunc = func(ctx context.Context, id int64) (*entity.User, error) {
					return &entity.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						FirstName: "Test",
						LastName:  "User",
					}, nil
				}
				m.updateFunc = func(ctx context.Context, user *entity.User) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "Error: Update fields are empty",
			userID: 1,
			requestBody: UpdateUserRequest{
				Username: nil,
				Email:    nil,
			},
			mockSetup: func(m *MockUserRepository) {
				m.getByIDFunc = func(ctx context.Context, id int64) (*entity.User, error) {
					return &entity.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						FirstName: "Test",
						LastName:  "User",
					}, nil
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			tt.mockSetup(mockRepo)

			handler := NewUserHandler(mockRepo)
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.Update(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectedError {
				var response entity.User
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if tt.requestBody.Username != nil && response.Username != *tt.requestBody.Username {
					t.Errorf("expected username %s, got %s", *tt.requestBody.Username, response.Username)
				}
				if tt.requestBody.Email != nil && response.Email != *tt.requestBody.Email {
					t.Errorf("expected email %s, got %s", *tt.requestBody.Email, response.Email)
				}
			}
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
