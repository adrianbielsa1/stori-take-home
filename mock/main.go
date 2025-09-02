package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := LoadConfig("config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("*path", echo)
	router.POST("*path", echo)

	if err := router.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatal(err)
	}
}

func echo(ctx *gin.Context) {
	var body map[string]any

	if ctx.Request.ContentLength > 0 {
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"headers": ctx.Request.Header,
		"query":   ctx.Request.URL.Query(),
		"body":    body,
	})
}
