package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/perlyanzagithub/property-service-common/services"
	"github.com/perlyanzagithub/property-service-common/utils"
	"net/http"
	"strings"
)

type JWTMiddleware struct {
	jwtService *services.JWTService
	redis      services.RedisService
}

func NewJWTMiddleware(jwtService *services.JWTService, redis services.RedisService) *JWTMiddleware {
	return &JWTMiddleware{jwtService: jwtService, redis: redis}
}

func (m *JWTMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.HandleResponse(c, http.StatusUnauthorized, "failed", nil, nil, nil)
			c.Abort()
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		decrypt, err := utils.DecryptAES(tokenStr)
		if err != nil {
			utils.HandleError(c, err, http.StatusUnauthorized)
			c.Abort()
		}
		if isBlacklisted, _ := m.redis.Get(decrypt); isBlacklisted != "active" {
			utils.HandleResponse(c, http.StatusUnauthorized, "failed", nil, nil, nil)
			c.Abort()
		}
		claims, err := m.jwtService.ParseToken(decrypt)
		if err != nil {
			utils.HandleError(c, err, http.StatusUnauthorized)
			c.Abort()
		}

		// Store user ID in context for further use
		for key, value := range claims {
			c.Set(key, value)
		}
		c.Next()
	}
}
