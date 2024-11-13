package v1

import (
	"github.com/GearFramework/GophKeeper/internal/config"
	"github.com/GearFramework/GophKeeper/internal/pkg/logger"
	"github.com/GearFramework/GophKeeper/internal/pkg/rest/middleware"
	"github.com/gin-gonic/gin"
)

// RestAPI struct of REST API
type RestAPI struct {
	conf   *config.Server
	router *gin.Engine
}

// NewRestAPI return new REST API
func NewRestAPI(conf *config.Server) *RestAPI {
	gin.SetMode(gin.ReleaseMode)
	return &RestAPI{
		conf:   conf,
		router: gin.New(),
	}
}

// GetRouter return router of API
func (a *RestAPI) GetRouter() any {
	return a.router
}

// Init REST API
func (a *RestAPI) Init() {
	logger.Log.Info("initializing REST API")
	logger.Log.Info("set middlewares")
	a.router.Use(middleware.Logger())
	a.router.Use(middleware.Compress())
	a.router.Use(middleware.Auth())
	a.initRoutes()
}

func (a *RestAPI) initRoutes() {
	logger.Log.Info("set routes")
	v1 := a.router.Group("/v1")
	a.addPing(v1)
	a.addSignup(v1)
	a.addSignin(v1)
	a.addEntity(v1)
	a.router.NoRoute(func(ctx *gin.Context) {
		logger.Log.Errorf("not found route %s", ctx.Request.URL)
		a.NotFoundHandler(ctx)
	})
}
