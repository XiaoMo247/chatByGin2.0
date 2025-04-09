package routers

import (
	"chatByGin5.0/controllers"

	"github.com/gin-gonic/gin"
)

func IndexRouterInit(c *gin.Engine) {
	IndexRouter := c.Group("/")
	{
		IndexRouter.GET("", controllers.IndexInit)
	}
}
