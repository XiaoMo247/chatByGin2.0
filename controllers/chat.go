package controllers

import (
	"fmt"
	"net/http"

	"chatByGin5.0/models"
	"chatByGin5.0/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ChatIndex(c *gin.Context) {
	// 从token中解析出user_id
	user_id, err := utils.ExtractTokenID(c)
	// fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 根据user_id从数据库查询数据
	u, err := models.GetUserByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "chat/index.html", gin.H{
		"Username": u.Username,
	})
}

func PublicChatIndex(c *gin.Context) {
	// 从token中解析出user_id
	user_id, err := utils.ExtractTokenID(c)
	// fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 根据user_id从数据库查询数据
	u, err := models.GetUserByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取历史消息
	historyMessages, err := models.GetHistoryMessages()
	if err != nil {
		fmt.Println("获取历史消息失败:", err)
		historyMessages = []map[string]string{}
	}

	c.HTML(http.StatusOK, "chat/public.html", gin.H{
		"Username":        u.Username,
		"HistoryMessages": historyMessages,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(hub *models.Hub, c *gin.Context) {
	// 从token中解析出user_id
	user_id, err := utils.ExtractTokenID(c)
	// fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 根据user_id从数据库查询数据
	u, err := models.GetUserByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 登录成功后检查并踢出旧会话
	for client := range hub.ClientList {
		if client.Username == u.Username {
			client.ForceLogout("重复登录")
			client.Logout()
		}
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade Error:", err)
	}

	newClient := models.NewClient(hub, ws)
	newClient.Username = u.Username

	hub.ChRegister <- newClient

	go newClient.Writer()
	go newClient.Reader()
}
