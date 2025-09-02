package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateRequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader("X-Request-ID") == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-Request-ID header"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ValidateAPIKey(configurationStore *ConfigurationStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config := configurationStore.Get()

		if ctx.GetHeader("x-api-key") != config.APIKey {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
