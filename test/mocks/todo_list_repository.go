package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/lumoshiveacademy/todolist/model"
	"github.com/stretchr/testify/mock"
)

// TodoListRepositoryMock is a testify mock for repository.TodoListRepository.
type TodoListRepositoryMock struct {
	mock.Mock
}

func (m *TodoListRepositoryMock) Create(ctx context.Context, todoList *model.TodoList) error {
	args := m.Called(ctx, todoList)
	return args.Error(0)
}

func (m *TodoListRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*model.TodoList, error) {
	args := m.Called(ctx, id)
	if val, ok := args.Get(0).(*model.TodoList); ok {
		return val, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *TodoListRepositoryMock) FindAll(ctx context.Context) ([]model.TodoList, error) {
	args := m.Called(ctx)
	if val, ok := args.Get(0).([]model.TodoList); ok {
		return val, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *TodoListRepositoryMock) Update(ctx context.Context, todoList *model.TodoList) error {
	args := m.Called(ctx, todoList)
	return args.Error(0)
}

func (m *TodoListRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
