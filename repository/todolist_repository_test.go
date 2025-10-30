package repository_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lumoshiveacademy/todolist/model"
	"github.com/lumoshiveacademy/todolist/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupRepository(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, repository.TodoListRepository, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	require.NoError(t, err)

	repo := repository.NewTodoListRepository(gormDB)

	cleanup := func() {
		require.NoError(t, sqlDB.Close())
	}

	return gormDB, mock, repo, cleanup
}

func TestTodoListRepository_Create_Succeeds(t *testing.T) {
	_, mock, repo, cleanup := setupRepository(t)
	defer func() {
		cleanup()
		require.NoError(t, mock.ExpectationsWereMet())
	}()

	todoList := &model.TodoList{
		ID:          uuid.New(),
		Title:       "Groceries",
		Description: "Weekly grocery items",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "todo_lists"`)).
		WithArgs(todoList.ID, todoList.Title, todoList.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectClose()

	err := repo.Create(context.Background(), todoList)
	require.NoError(t, err)
}

func TestTodoListRepository_FindByID_NotFound(t *testing.T) {
	_, mock, repo, cleanup := setupRepository(t)
	defer func() {
		cleanup()
		require.NoError(t, mock.ExpectationsWereMet())
	}()

	id := uuid.New()
	mock.ExpectQuery(`^SELECT \* FROM "todo_lists" WHERE id = \$1.*`).
		WithArgs(id, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "created_at", "updated_at"}))
	mock.ExpectClose()

	todoList, err := repo.FindByID(context.Background(), id)
	require.ErrorIs(t, err, repository.ErrTodoListNotFound)
	require.Nil(t, todoList)
}

func TestTodoListRepository_Delete_NotFound(t *testing.T) {
	_, mock, repo, cleanup := setupRepository(t)
	defer func() {
		cleanup()
		require.NoError(t, mock.ExpectationsWereMet())
	}()

	id := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec(`^DELETE FROM "todo_lists" WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mock.ExpectClose()

	err := repo.Delete(context.Background(), id)
	require.ErrorIs(t, err, repository.ErrTodoListNotFound)
}
