package main

import (
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

	docker.InitializeDockerClient()

	db.InitializeDatabase()

	controllers.InitializeRoutes()

	controllers.Router.Run(config.EnvConfig.Host + ":" + config.EnvConfig.Port)

	if err := db.CloseDatabase(); err != nil {
		panic("Failed to close database connection: " + err.Error())
	}
	if err := docker.StopAllContainers(); err != nil {
		panic("Failed to stop Docker containers: " + err.Error())
	}
	if err := docker.CloseDockerClient(); err != nil {
		panic("Failed to close Docker client: " + err.Error())
	}
	println("Server stopped gracefully")
}
