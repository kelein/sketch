package router

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	swagFile "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"sketch/pkg/middleware"
	"sketch/pkg/version"

	// Swagger Annotation Docs
	_ "sketch/pkg/swagger"
)

var app *gin.Engine

func init() { gin.SetMode(gin.ReleaseMode) }

func initApp() {
	app = gin.New()

	// Enable Logger and Recovery
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	// Enable Pprof Middleware
	prof := middleware.NewGinprof()
	prof.Register(app)

	// Enable Prometheus Middleware
	prom := middleware.NewProm(version.AppName)
	prom.Register(app)
}

func initRouter() {
	app.Any("/", index)
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swagFile.Handler))

	// Endpoints For v1
	v1 := app.Group("/v1")
	v1.Any("/", index)
}

// index godoc
// @Summary index
// @Description index
// @Tags index
// @ID /index
// @Accept json
// @Produce json
// @Success 200
func index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"app":    version.AppName,
		"pid":    os.Getegid(),
		"uptime": time.Now(),
		"build":  version.Info(),
	})
}

// Start ...
func Start() {
	initApp()
	initRouter()
	app.Run(":9000")
}
