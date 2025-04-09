package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Ip       string
	ChSend   chan []byte
	ChatData *ChatData
	Username string
}

func NewClient(hub *Hub, ws *websocket.Conn) *Client {
	return &Client{
		Hub:      hub,
		Conn:     ws,
		ChSend:   make(chan []byte, 256),
		ChatData: &ChatData{},
	}
}

var User_list = []string{}

func (c *Client) Writer() {
	defer func() {
		c.Logout()
	}()
	for message := range c.ChSend {
		c.Conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (c *Client) Reader() {
	defer func() {
		c.Logout()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		json.Unmarshal(message, &c.ChatData)
		switch c.ChatData.Type {
		case "login":
			c.Username = c.ChatData.Content
			c.ChatData.From = c.Username
			User_list = append(User_list, c.Username)
			c.ChatData.UserList = User_list
			LoginMessage, _ := json.Marshal(c.ChatData)
			c.Hub.ChBoardcast <- LoginMessage

		case "user":
			// 设置消息时间
			c.ChatData.Time = time.Now().Format("2006-01-02 15:04:05")
			// 保存消息到Redis
			err := SaveMessageToRedis(c.ChatData.From, c.ChatData.Content)
			if err != nil {
				fmt.Println("保存消息到Redis失败:", err)
			}
			c.ChatData.Type = "user"
			ChatMessage, _ := json.Marshal(c.ChatData)
			c.Hub.ChBoardcast <- ChatMessage
		}
	}
}

func (c *Client) Logout() {
	c.ChatData.Type = "logout"
	User_list = Del(User_list, c.Username)
	c.ChatData.UserList = User_list
	c.ChatData.Content = c.Username
	LogoutMessage, _ := json.Marshal(c.ChatData)
	c.Hub.ChBoardcast <- LogoutMessage
	c.Conn.Close()
	delete(c.Hub.ClientList, c)
}

func Del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return n_slice
}

func (c *Client) ForceLogout(reason string) {
	// 1. 发送被踢出通知消息
	kickoutMsg := &ChatData{
		Type:    "system",
		Content: "您的账号已在其他设备登录，您已被强制退出。原因: " + reason,
	}
	data, _ := json.Marshal(kickoutMsg)
	c.ChSend <- data

	// 2. 延迟关闭连接确保消息送达
	time.AfterFunc(100*time.Millisecond, func() {
		c.Conn.Close() // 关闭连接会触发客户端的onclose事件
	})
}
