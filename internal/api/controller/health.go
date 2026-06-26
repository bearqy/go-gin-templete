package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gin-templete/internal/api/util"
	"go-gin-templete/internal/home"
)

// Home
// @Tags Health
// @Summary 健康检测
// @Success 200 {object} util.APIResponse{data=string} "成功"
// @Router	/home [get]
func Home(c *gin.Context) {
	home.Home()
	c.JSON(http.StatusOK, util.ResponseSuccessful("Hello, World!"))
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, util.ResponseSuccessful(gin.H{"status": "ok"}))
}
