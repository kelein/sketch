package middleware

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

// GinProf pprof middleware for gin framework
type GinProf struct {
	Endpoint string
}

// NewGinProf return ginprof with prefix
func NewGinProf(prefixs ...string) *GinProf {
	prefix := "/debug/pprof"
	if len(prefixs) > 0 {
		prefix = prefixs[0]
	}
	return &GinProf{Endpoint: prefix}
}

// Register inject ginprof http routers
func (gf *GinProf) Register(app *gin.Engine) {
	gf.routers(app)
}

func (gf *GinProf) routers(app *gin.Engine) {
	g := app.Group(gf.Endpoint)

	g.GET("/", hwrap(pprof.Index))
	g.GET("/trace", hwrap(pprof.Trace))
	g.GET("/symbol", hwrap(pprof.Symbol))
	g.GET("/cmdline", hwrap(pprof.Cmdline))
	g.GET("/profile", hwrap(pprof.Profile))

	g.GET("/heap", hwrap(pprof.Handler("heap").ServeHTTP))
	g.GET("/block", hwrap(pprof.Handler("block").ServeHTTP))
	g.GET("/mutex", hwrap(pprof.Handler("mutex").ServeHTTP))
	g.GET("/allocs", hwrap(pprof.Handler("allocs").ServeHTTP))
	g.GET("/goroutine", hwrap(pprof.Handler("goroutine").ServeHTTP))
	g.GET("/threadcreate", hwrap(pprof.Handler("threadcreate").ServeHTTP))
}

func hwrap(h http.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		http.HandlerFunc(h).ServeHTTP(ctx.Writer, ctx.Request)
	}
}
