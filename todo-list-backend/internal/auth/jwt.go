package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joy_project/todo-list-backend/internal/models"
)

var jwtKey = []byte("your_secret_key") // 在生产环境中应该使用环境变量

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 为用户生成JWT令牌
func GenerateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT令牌并返回用户ID
func ValidateToken(tokenString string) (int, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
