package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"chatByGin5.0/models"
	"chatByGin5.0/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

func TokenValid(c *gin.Context) error {
	tokenString := utils.CookieGetToken(c)
	if tokenString == "" {
		return errors.New("未提供 Token")
	}

	// 解析 Token 获取 user_id
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}

	// 检查 Token 是否在数据库里有效
	var user models.User
	if err := models.DB.Where("token = ?", tokenString).First(&user).Error; err != nil {
		return errors.New("Token 已失效，请重新登录")
	}

	return nil
}
