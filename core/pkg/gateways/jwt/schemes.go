package jwt

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	ID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthContext struct {
	*gin.Context
}

type Token struct {
	Value     string
	ExpiredAt time.Time
}

func (c *AuthContext) GetUserID() (uint, error) {
	val, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user id not found in context")
	}

	id, ok := val.(uint)
	if !ok {
		return 0, errors.New("invalid user id type")
	}

	return id, nil
}
