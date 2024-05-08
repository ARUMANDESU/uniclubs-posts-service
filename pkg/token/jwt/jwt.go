package jwt

import (
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func GetUserIDFromToken(tokenString string, secret string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret used to sign the token
		return []byte(secret), nil
	})
	if err != nil {
		errorMessage := err.Error()
		switch {
		case strings.Contains(errorMessage, "token signature is invalid"):
			return 0, domain.ErrTokenSignatureIsInvalid
		case strings.Contains(errorMessage, "token is expired"):
			return 0, domain.ErrTokenIsExpired
		default:
			return 0, err
		}
	}

	// Check if the token is valid
	if !token.Valid {
		return 0, domain.ErrTokenIsNotValid
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, domain.ErrInvalidTokenClaims
	}

	// Extract user_id from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, domain.ErrUserIDClaimNotFound
	}

	return int64(userID), nil
}
