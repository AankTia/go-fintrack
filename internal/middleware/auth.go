package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/AankTia/go-fintrack/internal/platform"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(platform.AppConfig.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// JWT parse numbers sebagai float64, kita ubah ke string untuk DashboardHandler
		userID, ok := claims["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userIDStr := fmt.Sprintf("%v", userID) 
		// Jika userID adalah float64 (misal 123.0), fmt %v mungkin menghasilkan "123" atau "123.0"
		// Untuk lebih aman jika itu angka:
		if val, ok := userID.(float64); ok {
			userIDStr = fmt.Sprintf("%.0f", val)
		}

		ctx := context.WithValue(r.Context(), "user_id", userIDStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
