package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lumoshiveacademy/todolist/model"
	"github.com/lumoshiveacademy/todolist/repository"
	"go.uber.org/zap"
)

// TodoListService defines business operations for todo lists.
type TodoListService interface {
	CreateTodoList(ctx context.Context, req model.CreateTodoListRequest) (model.TodoListResponse, error)
	GetTodoList(ctx context.Context, id uuid.UUID) (model.TodoListResponse, error)
	ListTodoLists(ctx context.Context) ([]model.TodoListResponse, error)
	UpdateTodoList(ctx context.Context, id uuid.UUID, req model.UpdateTodoListRequest) (model.TodoListResponse, error)
	DeleteTodoList(ctx context.Context, id uuid.UUID) error
}

type todoListService struct {
	repository repository.TodoListRepository
	logger     *zap.Logger
}

// NewTodoListService constructs a TodoListService implementation.
func NewTodoListService(repository repository.TodoListRepository, logger *zap.Logger) TodoListService {
	return &todoListService{
		repository: repository,
		logger:     logger,
	}
}

func (s *todoListService) CreateTodoList(ctx context.Context, req model.CreateTodoListRequest) (model.TodoListResponse, error) {
	todoList := &model.TodoList{
		Title:       req.Title,
		Description: req.Description,
	}
	if err := s.repository.Create(ctx, todoList); err != nil {
		s.logger.Error("create todo list failed", zap.Error(err))
		return model.TodoListResponse{}, fmt.Errorf("create todo list: %w", err)
	}
	s.logger.Info("todo list created", zap.String("id", todoList.ID.String()))
	return todoList.ToResponse(), nil
}

func (s *todoListService) GetTodoList(ctx context.Context, id uuid.UUID) (model.TodoListResponse, error) {
	todoList, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrTodoListNotFound {
			return model.TodoListResponse{}, err
		}
		s.logger.Error("get todo list failed", zap.String("id", id.String()), zap.Error(err))
		return model.TodoListResponse{}, fmt.Errorf("get todo list: %w", err)
	}
	return todoList.ToResponse(), nil
}

func (s *todoListService) ListTodoLists(ctx context.Context) ([]model.TodoListResponse, error) {
	todoLists, err := s.repository.FindAll(ctx)
	if err != nil {
		s.logger.Error("list todo lists failed", zap.Error(err))
		return nil, fmt.Errorf("list todo lists: %w", err)
	}
	responses := make([]model.TodoListResponse, len(todoLists))
	for i, todoList := range todoLists {
		responses[i] = todoList.ToResponse()
	}
	return responses, nil
}

func (s *todoListService) UpdateTodoList(ctx context.Context, id uuid.UUID, req model.UpdateTodoListRequest) (model.TodoListResponse, error) {
	todoList, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrTodoListNotFound {
			return model.TodoListResponse{}, err
		}
		s.logger.Error("retrieve todo list for update failed", zap.String("id", id.String()), zap.Error(err))
		return model.TodoListResponse{}, fmt.Errorf("get todo list: %w", err)
	}

	todoList.Title = req.Title
	todoList.Description = req.Description

	if err := s.repository.Update(ctx, todoList); err != nil {
		s.logger.Error("update todo list failed", zap.String("id", id.String()), zap.Error(err))
		return model.TodoListResponse{}, fmt.Errorf("update todo list: %w", err)
	}
	s.logger.Info("todo list updated", zap.String("id", id.String()))
	return todoList.ToResponse(), nil
}

func (s *todoListService) DeleteTodoList(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		if err == repository.ErrTodoListNotFound {
			return err
		}
		s.logger.Error("delete todo list failed", zap.String("id", id.String()), zap.Error(err))
		return fmt.Errorf("delete todo list: %w", err)
	}
	s.logger.Info("todo list deleted", zap.String("id", id.String()))
	return nil
}
