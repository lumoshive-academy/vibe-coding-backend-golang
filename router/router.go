package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lumoshiveacademy/todolist/handler"
	appMiddleware "github.com/lumoshiveacademy/todolist/middleware"
	"go.uber.org/zap"
)

// New initializes the HTTP router with middleware and route registrations.
func New(todoListHandler *handler.TodoListHandler, logger *zap.Logger, jwtSecret, jwtIssuer string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(appMiddleware.Recovery(logger))
	r.Use(appMiddleware.Logger(logger))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(api chi.Router) {
		// api.Use(appMiddleware.JWTAuthentication(jwtSecret, jwtIssuer, logger))
		api.Route("/todolists", func(r chi.Router) {
			r.Post("/", todoListHandler.Create)
			r.Get("/", todoListHandler.List)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", todoListHandler.Get)
				r.Put("/", todoListHandler.Update)
				r.Delete("/", todoListHandler.Delete)
			})
		})
	})

	return r
}
