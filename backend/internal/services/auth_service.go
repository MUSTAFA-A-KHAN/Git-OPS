package services

import "github.com/golang-jwt/jwt/v5"

type AuthService struct{ secret string }

func NewAuthService(secret string) *AuthService { return &AuthService{secret: secret} }

func (a *AuthService) Token(userID string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": userID})
	return t.SignedString([]byte(a.secret))
}
