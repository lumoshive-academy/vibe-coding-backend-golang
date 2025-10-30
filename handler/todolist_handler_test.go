package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lumoshiveacademy/todolist/handler"
	"github.com/lumoshiveacademy/todolist/model"
	"github.com/lumoshiveacademy/todolist/package/response"
	"github.com/lumoshiveacademy/todolist/repository"
	"github.com/lumoshiveacademy/todolist/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestTodoListHandler_Create_Success(t *testing.T) {
	serviceMock := new(mocks.TodoListServiceMock)
	validate := validator.New(validator.WithRequiredStructEnabled())
	logger := zaptest.NewLogger(t)

	h := handler.NewTodoListHandler(serviceMock, validate, logger)

	reqBody := model.CreateTodoListRequest{Title: "Groceries", Description: "Weekly"}
	serviceMock.
		On("CreateTodoList", mock.Anything, reqBody).
		Return(model.TodoListResponse{
			ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Title:       "Groceries",
			Description: "Weekly",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil)

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/todolists", bytes.NewReader(body))

	h.Create(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var resp response.Message
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "success", resp.Status)
	serviceMock.AssertExpectations(t)
}

func TestTodoListHandler_Create_ValidationError(t *testing.T) {
	serviceMock := new(mocks.TodoListServiceMock)
	validate := validator.New(validator.WithRequiredStructEnabled())
	logger := zaptest.NewLogger(t)
	h := handler.NewTodoListHandler(serviceMock, validate, logger)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/todolists", bytes.NewReader([]byte(`{"description":"Weekly"}`)))

	h.Create(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	serviceMock.AssertNotCalled(t, "CreateTodoList", mock.Anything, mock.Anything)
}

func TestTodoListHandler_Get_NotFound(t *testing.T) {
	serviceMock := new(mocks.TodoListServiceMock)
	validate := validator.New(validator.WithRequiredStructEnabled())
	logger := zaptest.NewLogger(t)
	h := handler.NewTodoListHandler(serviceMock, validate, logger)

	id := uuid.New()
	serviceMock.
		On("GetTodoList", mock.Anything, id).
		Return(model.TodoListResponse{}, repository.ErrTodoListNotFound)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/todolists/"+id.String(), nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	h.Get(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	serviceMock.AssertExpectations(t)
}

func TestTodoListHandler_Delete_Success(t *testing.T) {
	serviceMock := new(mocks.TodoListServiceMock)
	validate := validator.New(validator.WithRequiredStructEnabled())
	logger := zaptest.NewLogger(t)
	h := handler.NewTodoListHandler(serviceMock, validate, logger)

	id := uuid.New()
	serviceMock.On("DeleteTodoList", mock.Anything, id).Return(nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/todolists/"+id.String(), nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", id.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	h.Delete(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
	serviceMock.AssertExpectations(t)
}
