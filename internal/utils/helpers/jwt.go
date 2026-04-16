package helpers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/savanyv/cartify/config"
)

const (
	jwtIssuer = "cartify-backend"
	jwtAudience = "cartify-api"
	jwtAccessExpiry = 24 * time.Hour
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Role string `json:"role"`
	TokenVersion int `json:"token_version"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateAccessToken(userID, username, email, role string, tokenVersion int) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type jwtService struct {
	secretKey []byte
}

func NewJWTService() JWTService {
	cfg := config.LoadConfig()
	secret := cfg.JWTSecret
	if secret == "" {
		panic("JWT_SECRET is not set")
	}

	return &jwtService{
		secretKey: []byte(secret),
	}
}

func (j *jwtService) GenerateAccessToken(userID, username, email, role string, tokenVersion int) (string, error) {
	if userID == "" || username == "" || email == "" || role == "" {
		return "", errors.New("invalid jwt payload")
	}

	claims := &JWTClaims{
		UserID: userID,
		Username: username,
		Email: email,
		Role: role,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: jwtIssuer,
			Audience: []string{jwtAudience},
			Subject: userID,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtAccessExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtService) GenerateRefreshToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("invalid user id")
	}

	claims := &jwt.RegisteredClaims{
		Issuer: jwtIssuer,
		Audience: []string{jwtAudience},
		Subject: userID,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKey
			}
			return j.secretKey, nil
		},
		jwt.WithIssuer(jwtIssuer),
		jwt.WithAudience(jwtAudience),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}