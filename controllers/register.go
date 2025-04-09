package controllers

import (
	"net/http"

	"chatByGin5.0/models"

	"github.com/gin-gonic/gin"
)

func RegisterInit(c *gin.Context) {
	c.HTML(http.StatusOK, "register/register.html", gin.H{})
}

func RegisterUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	user := &models.User{}

	// 用户名验证
	if len(username) < 4 || len(username) > 18 {
		c.JSON(http.StatusBadGateway, gin.H{
			"Message": "用户名需要大于 4 位且小于 18 位！",
		})
		return
	}
	err := models.DB.Model(models.User{}).Where("username = ?", username).Take(&user).Error
	if err == nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "用户名已存在！",
		})
		return
	}

	// 密码验证
	if len(password) < 6 || len(password) > 18 {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "密码需要大于 6 位且小于 18 位！",
		})
		return
	}
	if password != repassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "两次密码不一致！",
		})
		return
	}

	// 注册 + 自动JWT加密 --- BeforeSave()
	user = &models.User{Username: username, Password: password}
	_, err = user.SaveUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功！"})
}
