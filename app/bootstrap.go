package app

import (
	"fmt"
	"io"

	"github.com/eviltomorrow/aphrodite-web/cache"

	"github.com/gin-gonic/gin"

	"github.com/eviltomorrow/aphrodite-web/controller"
	"github.com/eviltomorrow/aphrodite-web/db"
	"github.com/eviltomorrow/aphrodite-web/middleware"
)

// Port port
var Port int = 8080

// PathHTML html
var PathHTML string

// Startup startup
func Startup() {
	db.BuildMySQL()
	cache.BuildRedis()

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(middleware.LogWriter())

	var engine = gin.New()
	engine.Use(middleware.LogInterceptor())
	engine.Use(gin.Recovery())

	engine.Static("/static", PathHTML)
	engine.LoadHTMLGlob(fmt.Sprintf("%s/*.html", PathHTML))
	engine.GET("/", controller.Index)

	engine.Run(fmt.Sprintf(":%d", Port))
}
