package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/kenwoo9y/todo-api-go/api/internal/entity"
	"github.com/kenwoo9y/todo-api-go/api/internal/repository"
	"github.com/kenwoo9y/todo-api-go/api/pkg/common"
)

type TaskHandler struct {
	repo repository.TaskRepository
}

func NewTaskHandler(repo repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Status      string     `json:"status"`
	OwnerID     int64      `json:"owner_id"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Status      *string    `json:"status,omitempty"`
	OwnerID     *int64     `json:"owner_id,omitempty"`
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/tasks":
		h.Create(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/tasks":
		h.GetAll(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/users/") && strings.HasSuffix(r.URL.Path, "/tasks"):
		h.GetByOwnerID(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/tasks/"):
		h.GetByID(w, r)
	case r.Method == http.MethodPatch && strings.HasPrefix(r.URL.Path, "/tasks/"):
		h.Update(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/tasks/"):
		h.Delete(w, r)
	default:
		common.ErrorJSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.HandleError(w, err)
		return
	}

	task := &entity.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      entity.TaskStatus(req.Status),
		OwnerID:     req.OwnerID,
	}

	if err := h.repo.Create(r.Context(), task); err != nil {
		common.HandleError(w, err)
		return
	}

	common.JSONResponse(w, http.StatusCreated, task)
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodGet) {
		return
	}

	tasks, err := h.repo.GetAll(r.Context())
	if err != nil {
		common.HandleError(w, err)
		return
	}

	common.JSONResponse(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodGet) {
		return
	}

	id, err := common.ExtractIDFromPath(r.URL.Path, "/tasks/")
	if err != nil {
		common.HandleError(w, common.ErrInvalidID)
		return
	}

	task, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		common.HandleError(w, err)
		return
	}

	if task == nil {
		common.HandleError(w, common.ErrNotFound)
		return
	}

	common.JSONResponse(w, http.StatusOK, task)
}

func (h *TaskHandler) GetByOwnerID(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodGet) {
		return
	}

	ownerID, err := common.ExtractOwnerIDFromPath(r.URL.Path)
	if err != nil {
		common.HandleError(w, err)
		return
	}

	tasks, err := h.repo.GetByOwnerID(r.Context(), ownerID)
	if err != nil {
		common.HandleError(w, err)
		return
	}

	common.JSONResponse(w, http.StatusOK, tasks)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodPatch) {
		return
	}

	id, err := common.ExtractIDFromPath(r.URL.Path, "/tasks/")
	if err != nil {
		common.HandleError(w, common.ErrInvalidID)
		return
	}

	existingTask, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		common.HandleError(w, err)
		return
	}

	if existingTask == nil {
		common.HandleError(w, common.ErrNotFound)
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.HandleError(w, err)
		return
	}

	if req.Title != nil {
		existingTask.Title = *req.Title
	}
	if req.Description != nil {
		existingTask.Description = *req.Description
	}
	if req.DueDate != nil {
		existingTask.DueDate = req.DueDate
	}
	if req.Status != nil {
		existingTask.Status = entity.TaskStatus(*req.Status)
	}
	if req.OwnerID != nil {
		existingTask.OwnerID = *req.OwnerID
	}

	if req.Title == nil && req.Description == nil && req.DueDate == nil && req.Status == nil && req.OwnerID == nil {
		common.ErrorJSONResponse(w, http.StatusBadRequest, "at least one field must be provided for update")
		return
	}

	if err := h.repo.Update(r.Context(), existingTask); err != nil {
		common.HandleError(w, err)
		return
	}

	common.JSONResponse(w, http.StatusOK, existingTask)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodDelete) {
		return
	}

	id, err := common.ExtractIDFromPath(r.URL.Path, "/tasks/")
	if err != nil {
		common.HandleError(w, common.ErrInvalidID)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		common.HandleError(w, err)
		return
	}

	common.JSONResponse(w, http.StatusNoContent, nil)
}
