package controllers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/config"
)

var Router = gin.Default()

func InitializeRoutes() {
	Router.Use(gin.Recovery())

	// corsConfig := cors.DefaultConfig()
	// Router.Use(cors.New(corsConfig))

	root := Router.Group("/")

	InitializeSpaRoutes(root)
	InitializeApiRoutes(root)
	InitializePocketBaseRoutes(root)
}

func InitializeApiRoutes(root *gin.RouterGroup) {
	api := root.Group("/api/v1")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	InitializeServerRoutes(api)
}

func InitializePocketBaseRoutes(root *gin.RouterGroup) {
	target, _ := url.Parse("http://" + config.EnvConfig.PbHost + ":" + config.EnvConfig.PbPort)
	proxy := httputil.NewSingleHostReverseProxy(target)

	origDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		origHost := r.Host
		origDirector(r)

		if r.Header.Get("X-Forwarded-Host") == "" {
			r.Header.Set("X-Forwarded-Host", origHost)
		}
		if r.Header.Get("X-Forwarded-Proto") == "" {
			if r.TLS != nil {
				r.Header.Set("X-Forwarded-Proto", "https")
			} else {
				r.Header.Set("X-Forwarded-Proto", "http")
			}
		}
	}

	handler := http.StripPrefix("/pb", proxy)

	g := root.Group("/pb")
	g.Any("", gin.WrapH(handler))
	g.Any("/*path", gin.WrapH(handler))
}
