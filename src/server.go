package src

import (
	"github.com/gin-gonic/gin"
	"github.com/moocss/apiserver/src/router"
	"github.com/moocss/apiserver/src/router/middleware"
)

// New returns a app instance
func New() *gin.Engine {
	// init db

	// Set gin mode.
	gin.SetMode(Conf.Core.Mode)

	// Create the Gin engine.
	g := gin.New()

	// Routes
	router.Load(
		// Cores
		g,
		// Middlwares
		middleware.VersionMiddleware(),
		// middleware.Logging(),
		middleware.RequestId(),
	)
	return g
}
