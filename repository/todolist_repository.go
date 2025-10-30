package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lumoshiveacademy/todolist/model"
	"gorm.io/gorm"
)

// ErrTodoListNotFound indicates that the todo list record does not exist.
var ErrTodoListNotFound = errors.New("todo list not found")

// TodoListRepository defines database operations for todo lists.
type TodoListRepository interface {
	Create(ctx context.Context, todoList *model.TodoList) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.TodoList, error)
	FindAll(ctx context.Context) ([]model.TodoList, error)
	Update(ctx context.Context, todoList *model.TodoList) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type todoListRepository struct {
	db *gorm.DB
}

// NewTodoListRepository constructs a TodoListRepository backed by GORM.
func NewTodoListRepository(db *gorm.DB) TodoListRepository {
	return &todoListRepository{db: db}
}

func (r *todoListRepository) Create(ctx context.Context, todoList *model.TodoList) error {
	if err := r.db.WithContext(ctx).Create(todoList).Error; err != nil {
		return fmt.Errorf("create todo list: %w", err)
	}
	return nil
}

func (r *todoListRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.TodoList, error) {
	var todoList model.TodoList
	if err := r.db.WithContext(ctx).First(&todoList, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTodoListNotFound
		}
		return nil, fmt.Errorf("find todo list: %w", err)
	}
	return &todoList, nil
}

func (r *todoListRepository) FindAll(ctx context.Context) ([]model.TodoList, error) {
	var todoLists []model.TodoList
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&todoLists).Error; err != nil {
		return nil, fmt.Errorf("find all todo lists: %w", err)
	}
	return todoLists, nil
}

func (r *todoListRepository) Update(ctx context.Context, todoList *model.TodoList) error {
	if err := r.db.WithContext(ctx).Save(todoList).Error; err != nil {
		return fmt.Errorf("update todo list: %w", err)
	}
	return nil
}

func (r *todoListRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.TodoList{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("delete todo list: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrTodoListNotFound
	}
	return nil
}
