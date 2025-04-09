package controllers

import (
	"net/http"

	"chatByGin5.0/models"
	"chatByGin5.0/utils"

	"github.com/gin-gonic/gin"
)

func LoginInit(c *gin.Context) {
	c.HTML(http.StatusOK, "login/login.html", gin.H{})
}

func PasswordForgot(c *gin.Context) {
	c.HTML(http.StatusOK, "login/pf.html", gin.H{})
}

func LoginUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	u := &models.User{
		Username: username,
		Password: password,
	}

	// 调用 models.LoginCheck 对用户名和密码进行验证
	token, err := models.LoginCheck(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username or password is incorrect.",
		})
		return
	}

	// 将 Token 存入 HttpOnly Cookie
	c.SetCookie("token", token, 3600, "/", "", true, true)

	c.Redirect(http.StatusFound, "/chat")
}

// 通过 token 获取用户
func CurrentUser(c *gin.Context) {
	// 从token中解析出user_id
	user_id, err := utils.ExtractTokenID(c)
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

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    u,
	})
}
