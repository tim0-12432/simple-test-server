package controllers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/config"
	"github.com/tim0-12432/simple-test-server/protocols"
)

var Router = gin.Default()

func InitializeRoutes() {
	Router.Use(gin.Recovery())

	allowedOrigins := []string{
		"http://" + config.EnvConfig.Host + ":" + config.EnvConfig.Port,
	}
	if config.EnvConfig.AllowedOrigins != nil {
		allowedOrigins = append(allowedOrigins, config.EnvConfig.AllowedOrigins...)
	}
	if config.EnvConfig.Env == "DEV" {
		allowedOrigins = append(allowedOrigins, "http://localhost:5173")
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = allowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	Router.Use(cors.New(corsConfig))

	root := Router.Group("/")

	InitializeSpaRoutes(root)
	InitializeApiRoutes(root)
	InitializePocketBaseRoutes(root)
}

func InitializeApiRoutes(root *gin.RouterGroup) {
	api := root.Group("/api/v1")

	InitializeServerRoutes(api)
	InitializeContainerRoutes(api)
	protocols.InitializeProtocolRoutes(api)
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
