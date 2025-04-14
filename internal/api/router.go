package api

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go-gin-templete/internal/api/controller"
	"net/http"
)

func Router(r gin.IRouter) {
	// Swagger 路由
	r.GET("/wiki", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/swagger/index.html") })
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/home", controller.Home)
	apiRouter := r.Group("/api")
	apiRouter.GET("/home", controller.Home)

}
