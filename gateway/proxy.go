package main

import (
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	Configuration *Configuration
}

func (p *Proxy) Handle(ctx *gin.Context) {
	service := ctx.Param("service")
	path := ctx.Param("path")

	target, ok := p.Configuration.Services[service]

	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Unknown service"})
		return
	}

	url := target + path

	fmt.Println(path)
	fmt.Println(url)

	req, err := http.NewRequest(ctx.Request.Method, url, ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}

	maps.Copy(req.Header, ctx.Request.Header)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("Error forwarding request: %w", err)})
		return
	}

	defer resp.Body.Close()

	maps.Copy(ctx.Writer.Header(), resp.Header)

	ctx.Status(resp.StatusCode)
	io.Copy(ctx.Writer, resp.Body)
}

func ReloadConfigHandler(store *ConfigurationStore, path string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newCfg, err := LoadConfig(path)

		if err != nil {
			log.Println("Config reload failed:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload config"})
			return
		}

		store.Update(newCfg)
		log.Println("Config reloaded successfully")
		ctx.JSON(http.StatusOK, gin.H{"status": "reloaded"})
	}
}
