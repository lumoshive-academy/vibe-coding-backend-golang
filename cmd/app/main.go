package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/lumoshiveacademy/todolist/database"
	"github.com/lumoshiveacademy/todolist/handler"
	"github.com/lumoshiveacademy/todolist/model"
	appConfig "github.com/lumoshiveacademy/todolist/package/config"
	appLogger "github.com/lumoshiveacademy/todolist/package/logger"
	"github.com/lumoshiveacademy/todolist/repository"
	"github.com/lumoshiveacademy/todolist/router"
	"github.com/lumoshiveacademy/todolist/service"
	"go.uber.org/zap"
)

func main() {
	cfg, err := appConfig.Load()
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	logger, err := appLogger.New(cfg.App.Name, cfg.App.Debug)
	if err != nil {
		panic(fmt.Errorf("init logger: %w", err))
	}
	defer func() { _ = logger.Sync() }()

	db, err := database.New(cfg.Database, logger)
	if err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}

	if err := db.AutoMigrate(&model.TodoList{}); err != nil {
		logger.Fatal("auto migrate failed", zap.Error(err))
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	todoListRepository := repository.NewTodoListRepository(db)
	todoListService := service.NewTodoListService(todoListRepository, logger)
	todoListHandler := handler.NewTodoListHandler(todoListService, validate, logger)

	httpRouter := router.New(todoListHandler, logger, cfg.JWT.Secret, cfg.JWT.Issuer)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      httpRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("starting http server", zap.Int("port", cfg.App.Port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("http server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", zap.Error(err))
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}

	logger.Info("server shutdown complete")
}
