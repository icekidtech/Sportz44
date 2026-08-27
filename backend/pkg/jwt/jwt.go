package jwt

import (
	"errors"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload for Sportz44.
type Claims struct {
	UserID uint   `json:"uid"`
	Role   string `json:"role"`
	gojwt.RegisteredClaims
}

// Generate issues an access/refresh token pair.
func Generate(userID uint, role, secret string, accessTTL, refreshTTL time.Duration) (access, refresh string, err error) {
	now := time.Now()

	accessClaims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: gojwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  gojwt.NewNumericDate(now),
			ExpiresAt: gojwt.NewNumericDate(now.Add(accessTTL)),
		},
	}
	access, err = gojwt.NewWithClaims(gojwt.SigningMethodHS256, accessClaims).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: gojwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  gojwt.NewNumericDate(now),
			ExpiresAt: gojwt.NewNumericDate(now.Add(refreshTTL)),
		},
	}
	refresh, err = gojwt.NewWithClaims(gojwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// Parse validates and decodes a token string.
func Parse(tokenStr, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := gojwt.ParseWithClaims(tokenStr, claims, func(t *gojwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
