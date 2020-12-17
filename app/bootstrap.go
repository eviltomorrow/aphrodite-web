package app

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/eviltomorrow/aphrodite-web/controller"
	"github.com/eviltomorrow/aphrodite-web/middleware"
)

// Port port
var Port int = 8080

// Startup startup
func Startup() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(middleware.LogWriter())

	var engine = gin.New()
	engine.Use(middleware.LogInterceptor())
	engine.Use(gin.Recovery())

	engine.GET("/", controller.Index)

	engine.Run(fmt.Sprintf(":%d", Port))
}
