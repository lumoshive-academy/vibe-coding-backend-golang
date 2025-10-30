package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lumoshiveacademy/todolist/package/response"
	"go.uber.org/zap"
)

type contextKey string

const (
	// ContextKeyClaims stores JWT claims in request context.
	ContextKeyClaims contextKey = "jwtClaims"
)

// JWTAuthentication validates JWT bearer tokens from the Authorization header.
func JWTAuthentication(secret, issuer string, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				unauthorized(w)
				return
			}

			tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || tokenString == "" {
				unauthorized(w)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					return nil, jwt.ErrTokenSignatureInvalid
				}
				return []byte(secret), nil
			}, jwt.WithAudience(issuer), jwt.WithIssuer(issuer))

			if err != nil || !token.Valid {
				logger.Warn("jwt validation failed", zap.Error(err))
				unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyClaims, token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	response.Write(w, http.StatusUnauthorized, response.Failure(map[string]string{
		"message": "unauthorized",
	}))
}
