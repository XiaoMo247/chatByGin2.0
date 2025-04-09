package routers

import (
	"chatByGin5.0/controllers"

	"github.com/gin-gonic/gin"
)

func LoginRouterInit(r *gin.Engine) {
	LoginRouter := r.Group("/login")
	{
		LoginRouter.GET("", controllers.LoginInit)
		LoginRouter.POST("", controllers.LoginUser)

		LoginRouter.GET("/pf", controllers.PasswordForgot)
	}
}
