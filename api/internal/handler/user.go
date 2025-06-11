package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kenwoo9y/todo-api-go/api/internal/entity"
	"github.com/kenwoo9y/todo-api-go/api/internal/repository"
	"github.com/kenwoo9y/todo-api-go/api/pkg/common"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateUserRequest struct {
	Username  *string `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/users":
		h.Create(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/users":
		h.GetAll(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/users/username/"):
		h.GetByUsername(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/users/") && !strings.Contains(r.URL.Path, "/tasks"):
		h.GetByID(w, r)
	case r.Method == http.MethodPatch && strings.HasPrefix(r.URL.Path, "/users/"):
		h.Update(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/users/"):
		h.Delete(w, r)
	default:
		common.ErrorJSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.HandleError(w, err)
		return
	}

	user := &entity.User{
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.repo.Create(r.Context(), user); err != nil {
		common.HandleError(w, err)
		return
	}

	common.JSONResponse(w, http.StatusCreated, user)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodGet) {
		return
	}

	users, err := h.repo.GetAll(r.Context())
	if err != nil {
		common.HandleError(w, err)
		return
	}

	common.JSONResponse(w, http.StatusOK, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodGet) {
		return
	}

	id, err := common.ExtractIDFromPath(r.URL.Path, "/users/")
	if err != nil {
		common.HandleError(w, common.ErrInvalidID)
		return
	}

	user, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		common.HandleError(w, err)
		return
	}

	if user == nil {
		common.HandleError(w, common.ErrNotFound)
		return
	}

	common.JSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodGet) {
		return
	}

	username := strings.TrimPrefix(r.URL.Path, "/users/username/")
	if username == "" {
		common.HandleError(w, common.ErrInvalidID)
		return
	}

	user, err := h.repo.GetByUsername(r.Context(), username)
	if err != nil {
		common.HandleError(w, err)
		return
	}

	if user == nil {
		common.HandleError(w, common.ErrNotFound)
		return
	}

	common.JSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodPatch) {
		return
	}

	id, err := common.ExtractIDFromPath(r.URL.Path, "/users/")
	if err != nil {
		common.HandleError(w, common.ErrInvalidID)
		return
	}

	existingUser, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		common.HandleError(w, err)
		return
	}

	if existingUser == nil {
		common.HandleError(w, common.ErrNotFound)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.HandleError(w, err)
		return
	}

	if req.Username != nil {
		existingUser.Username = *req.Username
	}
	if req.Email != nil {
		existingUser.Email = *req.Email
	}
	if req.FirstName != nil {
		existingUser.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		existingUser.LastName = *req.LastName
	}

	if req.Username == nil && req.Email == nil && req.FirstName == nil && req.LastName == nil {
		common.ErrorJSONResponse(w, http.StatusBadRequest, "at least one field must be provided for update")
		return
	}

	if err := h.repo.Update(r.Context(), existingUser); err != nil {
		common.HandleError(w, err)
		return
	}

	common.JSONResponse(w, http.StatusOK, existingUser)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !common.ValidateRequestMethod(w, r, http.MethodDelete) {
		return
	}

	id, err := common.ExtractIDFromPath(r.URL.Path, "/users/")
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
