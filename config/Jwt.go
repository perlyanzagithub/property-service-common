package config

import (
	"os"
	"time"
)

type JWTConfig struct {
	SecretKey      string
	ExpirationTime time.Duration
}

func LoadJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey:      os.Getenv("JWT_SECRET_KEY"),
		ExpirationTime: time.Hour * 24, // or load from env, e.g., os.Getenv("JWT_EXPIRATION"),
	}
}
