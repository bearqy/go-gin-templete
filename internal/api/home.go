package api

import (
	"net/http"

	"github.com/bearqy/go-gin-templete/internal/home"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	home.Home()
	c.String(http.StatusOK, "Hello, World!")
}
