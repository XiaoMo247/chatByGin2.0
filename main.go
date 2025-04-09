package main

import (
	"log"

	"chatByGin5.0/config"
	"chatByGin5.0/models"
	"chatByGin5.0/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConfigInit()
	models.MySQLInit()
	models.RedisInit()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("views/**/*")

	hub := models.NewHub()
	go hub.Run()

	routers.IndexRouterInit(r)
	routers.RegisterRouterInit(r)
	routers.LoginRouterInit(r)
	routers.ChatRouterInit(hub, r)

	r.Run(":1111")
}
