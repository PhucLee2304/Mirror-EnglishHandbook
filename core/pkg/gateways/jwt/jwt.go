package jwt

import (
	"core/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(claims *TokenClaims, secret []byte, ttl string) (string, error) {
	ttlDuration, err := time.ParseDuration(ttl)
	if err != nil {
		return "", err
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(ttlDuration))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateAccessToken(cfg *config.Config, id uint) (*Token, error) {
	claims := &TokenClaims{
		ID: id,
	}

	accessToken, err := createJWT(claims, []byte(cfg.AccessSecret), cfg.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	return &Token{
		Value:     accessToken,
		ExpiredAt: claims.ExpiresAt.Time,
	}, nil
}

func GenerateRefreshToken(cfg *config.Config, id uint) (*Token, error) {
	claims := &TokenClaims{
		ID: id,
	}

	refreshToken, err := createJWT(claims, []byte(cfg.RefreshSecret), cfg.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &Token{
		Value:     refreshToken,
		ExpiredAt: claims.ExpiresAt.Time,
	}, nil
}
