package router

import (
	"github.com/moocss/apiserver/src/api/sd"
	"github.com/moocss/apiserver/src/pkg/version"
	"github.com/moocss/apiserver/src/router/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moocss/apiserver/src/api/user"
)

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to api server.",
	})
}

func versionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"source":  "https://github.com/go-impatient/apiserver",
		"version": version.GetVersion(),
	})
}

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)

	g.GET("/version", versionHandler)
	g.GET("/", rootHandler)

	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "不存在的接口地址.")
	})

	// User API
	u := g.Group("/v1/user")
	{
		// u.POST("", user.Create)
		u.POST("", user.Create)
	}

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
