package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/perlyanzagithub/property-service-common/config"
	"os"
	"time"
)

type TokenUtil struct {
	config config.JWTConfig
}

func NewTokenUtil(cfg config.JWTConfig) *TokenUtil {
	return &TokenUtil{
		config: cfg,
	}
}

// GenerateToken creates a JWT token and stores it in Redis
func (t *TokenUtil) GenerateToken(data map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	for key, value := range data {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	// Store token in Redis
	//err = t.redisService.Set(tokenStr, "active", t.config.ExpirationTime)
	//if err != nil {
	//	return "", fmt.Errorf("failed to store token in Redis: %v", err)
	//}

}

func (t *TokenUtil) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	// Check if token exists in Redis
	//_, err := t.redisService.Get(tokenStr)
	//if err == redis.Nil {
	//	return nil, fmt.Errorf("token not found or expired")
	//} else if err != nil {
	//	return nil, err
	//}

	// Parse and validate token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.NewValidationError("invalid token claims", jwt.ValidationErrorClaimsInvalid)
	}

	return claims, nil
}
