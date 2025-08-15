package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var distPath = "./frontend/dist"

func InitializeSpaRoutes(root *gin.RouterGroup) {
	spa := root.Group("/spa")

	root.GET("/", func(c *gin.Context) {
		c.File(distPath + "/index.html")
	})

	spa.StaticFS("/", http.Dir(distPath+"/static"))

}
