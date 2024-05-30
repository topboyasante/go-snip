package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/topboyasante/go-snip/pkg/config"
)

func CreateJWTTokens(data any) (string, string, error) {
	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": data,                                       // Subject (user identifier)
		"iss": "invxice",                                  // Issuer
		"aud": data,                                       // Audience (user role)
		"exp": time.Now().Add(time.Hour * 24 * 1).Unix(), // Expiration time = 1 day
		"iat": time.Now().Unix(),                          // Issued at
	})

	accessTokenString, err := accessTokenClaims.SignedString([]byte(config.ENV.JWTKey))
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": data,                                       // Subject (user identifier)
		"iss": "invxice",                                  // Issuer
		"aud": data,                                       // Audience (user role)
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Expiration time = 30 days
		"iat": time.Now().Unix(),                          // Issued at
	})

	refreshTokenString, err := refreshTokenClaims.SignedString([]byte(config.ENV.JWTKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
