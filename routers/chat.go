package routers

import (
	"chatByGin5.0/controllers"
	"chatByGin5.0/middlewares"
	"chatByGin5.0/models"

	"github.com/gin-gonic/gin"
)

func ChatRouterInit(hub *models.Hub, c *gin.Engine) {
	ChatRouter := c.Group("/chat")
	{
		ChatRouter.Use(middlewares.JwtAuthMiddleware())
		ChatRouter.GET("", controllers.ChatIndex)

		ChatRouter.GET("/public", controllers.PublicChatIndex)
		ChatRouter.GET("/public/ws", func(c *gin.Context) {
			controllers.HandleWebSocket(hub, c)
		})
	}
}
