package web

import (
	"github.com/gin-gonic/gin"

	"github.com/LiangXianSen/gin-demo/config"
	M "github.com/LiangXianSen/go-utils/middleware"
)

type RouteGroup struct {
	BasePath    string
	Middlewares []gin.HandlerFunc
	Routes      []Route
}

type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

// LoadRouter loads routes returns router.
func LoadRouter(routeGroups []RouteGroup, conf *config.Config) *gin.Engine {
	routeGroups = InitMiddleware(routeGroups, conf)
	router := gin.New()
	for _, rg := range routeGroups {
		group := router.Group(rg.BasePath)
		group.Use(rg.Middlewares...)
		for _, route := range rg.Routes {
			group.Handle(route.Method, route.Path, route.Handler)
		}
	}
	return router
}

// InitMiddleware preparing work before running.
func InitMiddleware(routeGroups []RouteGroup, conf *config.Config) []RouteGroup {
	// logger setting
	loggerOpts := M.LoggerOptions{
		Application:  conf.General.Name,
		Version:      conf.General.Version,
		EnableOutput: true,
		EnableDebug:  conf.General.Debug,
	}

	// gin context injector
	// sets values into ctx of gin
	injector := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("conf", conf)
		}
	}()

	var groups []RouteGroup
	for _, rg := range routeGroups {
		rg.Middlewares = append([]gin.HandlerFunc{M.LoggerM(loggerOpts)}, rg.Middlewares...)
		rg.Middlewares = append(rg.Middlewares, injector)
		groups = append(groups, rg)
	}

	return groups
}
