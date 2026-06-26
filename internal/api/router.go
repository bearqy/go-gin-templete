package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "go-gin-templete/docs"
	"go-gin-templete/internal/api/controller"
)

func Router(r gin.IRouter) {
	// Swagger 路由
	r.GET("/wiki", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/swagger/index.html") })
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Backward-compatible examples.
	r.GET("/home", controller.Home)
	apiRouter := r.Group("/api")
	apiRouter.GET("/home", controller.Home)

	v1 := apiRouter.Group("/v1")
	v1.GET("/home", controller.Home)
	v1.GET("/health", controller.Health)
}
