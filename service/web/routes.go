package web

import (
	"github.com/gin-gonic/gin"

	M "github.com/LiangXianSen/go-utils/middleware"
)

/* Register routes here */
var routeGroups = []RouteGroup{
	{
		BasePath: "/",
		Middlewares: []gin.HandlerFunc{
			gin.Recovery(),
			M.RequestID,
			M.FaultHandler,
		},
		Routes: []Route{
			{"GET", "/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }},
		},
	},
}
