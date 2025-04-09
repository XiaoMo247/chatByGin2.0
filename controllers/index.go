package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexInit(c *gin.Context) {
	c.HTML(http.StatusOK, "index/index.html", gin.H{})
}
