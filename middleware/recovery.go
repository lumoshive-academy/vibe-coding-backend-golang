package middleware

import (
	"net/http"

	"github.com/lumoshiveacademy/todolist/package/response"
	"go.uber.org/zap"
)

// Recovery recovers from panics and writes a JSON error response.
func Recovery(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered", zap.Any("error", rec))
					response.Write(w, http.StatusInternalServerError, response.Failure(map[string]string{
						"message": "internal server error",
					}))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
