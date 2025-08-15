package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/config"
	"github.com/tim0-12432/simple-test-server/controllers"
	"github.com/tim0-12432/simple-test-server/db"
	"github.com/tim0-12432/simple-test-server/docker"
)

func main() {
	config.InitializeEnvConfig()

	if config.EnvConfig.Env == "DEV" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	db.InitializeDatabase()

	controllers.InitializeRoutes()

	srv := &http.Server{
		Addr:    config.EnvConfig.Host + ":" + config.EnvConfig.Port,
		Handler: controllers.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v", err)
	}

	if err := db.CloseDatabase(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
	}
	if err := docker.StopAllContainers(); err != nil {
		log.Printf("Failed to stop Docker containers: %v", err)
	}

	log.Println("Server stopped gracefully")
}
