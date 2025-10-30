package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/lumoshiveacademy/todolist/model"
	"github.com/stretchr/testify/mock"
)

// TodoListServiceMock is a testify mock for service.TodoListService.
type TodoListServiceMock struct {
	mock.Mock
}

func (m *TodoListServiceMock) CreateTodoList(ctx context.Context, req model.CreateTodoListRequest) (model.TodoListResponse, error) {
	args := m.Called(ctx, req)
	if resp, ok := args.Get(0).(model.TodoListResponse); ok {
		return resp, args.Error(1)
	}
	return model.TodoListResponse{}, args.Error(1)
}

func (m *TodoListServiceMock) GetTodoList(ctx context.Context, id uuid.UUID) (model.TodoListResponse, error) {
	args := m.Called(ctx, id)
	if resp, ok := args.Get(0).(model.TodoListResponse); ok {
		return resp, args.Error(1)
	}
	return model.TodoListResponse{}, args.Error(1)
}

func (m *TodoListServiceMock) ListTodoLists(ctx context.Context) ([]model.TodoListResponse, error) {
	args := m.Called(ctx)
	if resp, ok := args.Get(0).([]model.TodoListResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *TodoListServiceMock) UpdateTodoList(ctx context.Context, id uuid.UUID, req model.UpdateTodoListRequest) (model.TodoListResponse, error) {
	args := m.Called(ctx, id, req)
	if resp, ok := args.Get(0).(model.TodoListResponse); ok {
		return resp, args.Error(1)
	}
	return model.TodoListResponse{}, args.Error(1)
}

func (m *TodoListServiceMock) DeleteTodoList(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
