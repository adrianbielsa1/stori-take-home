package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

var store *ConfigurationStore

func main() {
	cfg, err := LoadConfig("config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	store = NewConfigStore(cfg)

	go func() {
		for {
			time.Sleep(30 * time.Second)

			newCfg, err := LoadConfig("config.yaml")

			if err != nil {
				log.Println("Error reloading config:", err)
				continue
			}

			store.Update(newCfg)
		}
	}()

	proxy := &Proxy{Configuration: store.Get()}

	router := gin.Default()

	// Middlewares.
	router.Use(ValidateRequestID())
	router.Use(ValidateAPIKey(store))

	router.GET("/api/:service/*path", proxy.Handle)
	router.POST("/api/:service/*path", proxy.Handle)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
