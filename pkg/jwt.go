package pkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type customClaims struct {
	jwt.RegisteredClaims
}

// ValidateJwtAndExtractSubject проверяет JWT токен и извлекает из него subject
func ValidateJwtAndExtractSubject(raw string, hmacsecret string) (string, error) {
	// https://github.com/dgrijalva/jwt-go/pull/437
	token, err := jwt.ParseWithClaims(raw, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// почему ВАЖНО проверять алгоритм подписи токена
			// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
			return nil, fmt.Errorf("неожиданный алгоритм, используемый для подписи: %v", token.Header["alg"])
		}
		return []byte(hmacsecret), nil
	})

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	issuedAt, err := token.Claims.GetIssuedAt()
	if err != nil {
		return "", err
	}
	if issuedAt.After(time.Now()) {
		return "", fmt.Errorf("token issued in future")
	}
	expireAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", err
	}
	if expireAt.Before(time.Now()) {
		return "", fmt.Errorf("token expired")
	}
	subj, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	return subj, nil
}
