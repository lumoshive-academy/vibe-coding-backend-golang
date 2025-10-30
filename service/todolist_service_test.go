package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lumoshiveacademy/todolist/model"
	"github.com/lumoshiveacademy/todolist/repository"
	"github.com/lumoshiveacademy/todolist/service"
	"github.com/lumoshiveacademy/todolist/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestTodoListService_CreateTodoList_Success(t *testing.T) {
	mockRepo := new(mocks.TodoListRepositoryMock)
	logger := zaptest.NewLogger(t)
	svc := service.NewTodoListService(mockRepo, logger)

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(todoList *model.TodoList) bool {
		return todoList.Title == "Groceries"
	})).Return(nil).Run(func(args mock.Arguments) {
		todoList := args.Get(1).(*model.TodoList)
		todoList.ID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
		todoList.CreatedAt = time.Now()
		todoList.UpdatedAt = time.Now()
	})

	res, err := svc.CreateTodoList(context.Background(), model.CreateTodoListRequest{Title: "Groceries"})
	require.NoError(t, err)
	require.Equal(t, "Groceries", res.Title)
	require.Equal(t, uuid.MustParse("11111111-1111-1111-1111-111111111111"), res.ID)
	mockRepo.AssertExpectations(t)
}

func TestTodoListService_GetTodoList_NotFound(t *testing.T) {
	mockRepo := new(mocks.TodoListRepositoryMock)
	logger := zaptest.NewLogger(t)
	svc := service.NewTodoListService(mockRepo, logger)

	id := uuid.New()
	mockRepo.On("FindByID", mock.Anything, id).Return(nil, repository.ErrTodoListNotFound)

	_, err := svc.GetTodoList(context.Background(), id)
	require.ErrorIs(t, err, repository.ErrTodoListNotFound)
	mockRepo.AssertExpectations(t)
}

func TestTodoListService_UpdateTodoList_Success(t *testing.T) {
	mockRepo := new(mocks.TodoListRepositoryMock)
	logger := zaptest.NewLogger(t)
	svc := service.NewTodoListService(mockRepo, logger)

	id := uuid.New()
	existing := &model.TodoList{ID: id, Title: "Old", Description: "old"}

	mockRepo.On("FindByID", mock.Anything, id).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(todoList *model.TodoList) bool {
		return todoList.Title == "New"
	})).Return(nil)

	res, err := svc.UpdateTodoList(context.Background(), id, model.UpdateTodoListRequest{Title: "New"})
	require.NoError(t, err)
	require.Equal(t, "New", res.Title)
	mockRepo.AssertExpectations(t)
}

func TestTodoListService_DeleteTodoList_Error(t *testing.T) {
	mockRepo := new(mocks.TodoListRepositoryMock)
	logger := zaptest.NewLogger(t)
	svc := service.NewTodoListService(mockRepo, logger)

	id := uuid.New()
	mockRepo.On("Delete", mock.Anything, id).Return(repository.ErrTodoListNotFound)

	err := svc.DeleteTodoList(context.Background(), id)
	require.ErrorIs(t, err, repository.ErrTodoListNotFound)
	mockRepo.AssertExpectations(t)
}
