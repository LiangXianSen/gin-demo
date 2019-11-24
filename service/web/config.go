package web

import (
	"github.com/gin-gonic/gin"

	M "github.com/LiangXianSen/gin-demo/middleware"
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

func LoadRouter(routeGroups []RouteGroup) *gin.Engine {
	routeGroups = InitMiddleware(routeGroups)
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

func InitMiddleware(routeGroups []RouteGroup) []RouteGroup {
	// logger setting
	loggerOpts := M.LoggerOptions{
		Application:  "demo",
		Version:      "v0.0.1",
		EnableOutput: true,
		EnableDebug:  true,
	}

	// gin context injector
	// sets values into ctx of gin
	injector := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("conf", "")
		}
	}()

	var groups []RouteGroup
	for _, rg := range routeGroups {
		rg.Middlewares = append(rg.Middlewares, injector)
		rg.Middlewares = append(rg.Middlewares, M.LoggerM(loggerOpts))
		groups = append(groups, rg)
	}

	return groups
}
