package routers

import (
	"chatByGin5.0/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRouterInit(r *gin.Engine) {
	RegisterRouter := r.Group("/register")
	{
		RegisterRouter.GET("", controllers.RegisterInit)

		RegisterRouter.POST("", controllers.RegisterUser)
	}
}
