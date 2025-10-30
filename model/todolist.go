package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TodoList represents a collection of todo items owned by the user.
type TodoList struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title       string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// BeforeCreate ensures the TodoList has a UUID before persisting.
func (t *TodoList) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// CreateTodoListRequest defines the expected payload for creating a todo list.
type CreateTodoListRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=1024"`
}

// UpdateTodoListRequest defines the payload for updating a todo list.
type UpdateTodoListRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"max=1024"`
}

// TodoListResponse describes the response returned to clients.
type TodoListResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts the model into a response DTO.
func (t TodoList) ToResponse() TodoListResponse {
	return TodoListResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
