// service.go — signs and validates JWTs.
// In prod you'd check credentials against DB — here we use a hardcoded user for practice.
package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return []byte("dev-secret-change-in-prod")
	}
	return []byte(s)
}

func Login(input LoginInput) (*TokenResponse, error) {
	if input.Email != "admin@test.com" || input.Password != "password" {
		return nil, errors.New("invalid credentials")
	}

	claims := &Claims{
		UserID: 1,
		Email:  input.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret())
	if err != nil {
		return nil, err
	}
	return &TokenResponse{Token: signed}, nil
}

// ValidateToken parses and validates a JWT string, returns claims if valid.
func ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(t *jwt.Token) (interface{}, error) {
			// Verify algorithm matches what we expect — prevents algorithm confusion attacks
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret(), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
