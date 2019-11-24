package web

import (
	"github.com/gin-gonic/gin"

	M "github.com/LiangXianSen/gin-demo/middleware"
)

/* Register routes here */
var routeGroups = []RouteGroup{
	{
		BasePath: "/",
		Middlewares: []gin.HandlerFunc{
			gin.Recovery(),
			gin.Logger(),
			M.RequestID,
			//M.LoginRequired,
			M.FaultHandler,
		},
		Routes: []Route{
			{"GET", "/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }},
		},
	},
}
