package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lumoshiveacademy/todolist/model"
	"github.com/lumoshiveacademy/todolist/package/response"
	"github.com/lumoshiveacademy/todolist/repository"
	"github.com/lumoshiveacademy/todolist/service"
	"go.uber.org/zap"
)

// TodoListHandler exposes HTTP handlers for todo list resources.
type TodoListHandler struct {
	service  service.TodoListService
	validate *validator.Validate
	logger   *zap.Logger
}

// NewTodoListHandler constructs a TodoListHandler.
func NewTodoListHandler(service service.TodoListService, validate *validator.Validate, logger *zap.Logger) *TodoListHandler {
	return &TodoListHandler{
		service:  service,
		validate: validate,
		logger:   logger,
	}
}

// Create handles POST /todolists requests.
func (h *TodoListHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTodoListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid todo list create payload", zap.Error(err))
		response.Write(w, http.StatusBadRequest, response.Failure(map[string]string{
			"message": "invalid request payload",
		}))
		return
	}

	if err := h.validate.StructCtx(r.Context(), req); err != nil {
		h.logger.Warn("todo list create validation failed", zap.Error(err))
		response.Write(w, http.StatusUnprocessableEntity, response.Failure(validationErrors(err)))
		return
	}

	todoList, err := h.service.CreateTodoList(r.Context(), req)
	if err != nil {
		h.logger.Error("todo list creation failed", zap.Error(err))
		response.Write(w, http.StatusInternalServerError, response.Failure(map[string]string{
			"message": "could not create todo list",
		}))
		return
	}

	response.Write(w, http.StatusCreated, response.Success(todoList))
}

// List handles GET /todolists requests.
func (h *TodoListHandler) List(w http.ResponseWriter, r *http.Request) {
	todoLists, err := h.service.ListTodoLists(r.Context())
	if err != nil {
		h.logger.Error("list todo lists failed", zap.Error(err))
		response.Write(w, http.StatusInternalServerError, response.Failure(map[string]string{
			"message": "could not fetch todo lists",
		}))
		return
	}

	response.Write(w, http.StatusOK, response.Success(todoLists))
}

// Get handles GET /todolists/{id} requests.
func (h *TodoListHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Write(w, http.StatusBadRequest, response.Failure(map[string]string{
			"message": "invalid todo list id",
		}))
		return
	}

	todoList, err := h.service.GetTodoList(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTodoListNotFound) {
			response.Write(w, http.StatusNotFound, response.Failure(map[string]string{
				"message": "todo list not found",
			}))
			return
		}
		h.logger.Error("get todo list failed", zap.String("id", id.String()), zap.Error(err))
		response.Write(w, http.StatusInternalServerError, response.Failure(map[string]string{
			"message": "could not fetch todo list",
		}))
		return
	}

	response.Write(w, http.StatusOK, response.Success(todoList))
}

// Update handles PUT /todolists/{id} requests.
func (h *TodoListHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Write(w, http.StatusBadRequest, response.Failure(map[string]string{
			"message": "invalid todo list id",
		}))
		return
	}

	var req model.UpdateTodoListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid todo list update payload", zap.Error(err))
		response.Write(w, http.StatusBadRequest, response.Failure(map[string]string{
			"message": "invalid request payload",
		}))
		return
	}

	if err := h.validate.StructCtx(r.Context(), req); err != nil {
		h.logger.Warn("todo list update validation failed", zap.Error(err))
		response.Write(w, http.StatusUnprocessableEntity, response.Failure(validationErrors(err)))
		return
	}

	todoList, err := h.service.UpdateTodoList(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrTodoListNotFound) {
			response.Write(w, http.StatusNotFound, response.Failure(map[string]string{
				"message": "todo list not found",
			}))
			return
		}
		h.logger.Error("update todo list failed", zap.String("id", id.String()), zap.Error(err))
		response.Write(w, http.StatusInternalServerError, response.Failure(map[string]string{
			"message": "could not update todo list",
		}))
		return
	}

	response.Write(w, http.StatusOK, response.Success(todoList))
}

// Delete handles DELETE /todolists/{id} requests.
func (h *TodoListHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Write(w, http.StatusBadRequest, response.Failure(map[string]string{
			"message": "invalid todo list id",
		}))
		return
	}

	if err := h.service.DeleteTodoList(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrTodoListNotFound) {
			response.Write(w, http.StatusNotFound, response.Failure(map[string]string{
				"message": "todo list not found",
			}))
			return
		}
		h.logger.Error("delete todo list failed", zap.String("id", id.String()), zap.Error(err))
		response.Write(w, http.StatusInternalServerError, response.Failure(map[string]string{
			"message": "could not delete todo list",
		}))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	id := chi.URLParam(r, param)
	return uuid.Parse(id)
}

func validationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range validationErrs {
			errorsMap[strings.ToLower(validationErr.Field())] = validationErr.Error()
		}
		return errorsMap
	}
	errorsMap["message"] = err.Error()
	return errorsMap
}
