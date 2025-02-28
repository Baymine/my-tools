package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/joy_project/todo-list-backend/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"

// Auth 中间件验证JWT令牌并将用户ID添加到请求上下文中
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从Authorization头获取令牌
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 检查Bearer前缀
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header must be in the format 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		// 验证令牌
		userID, err := auth.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 将用户ID添加到请求上下文
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserID 从请求上下文中获取用户ID
func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}
