package api

import (
	_ "github.com/cmrd-a/gophermart/internal/api/docs"
	"github.com/cmrd-a/gophermart/internal/api/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Gophermart API
//	@version		1.0
//	@description	Накопительная система лояльности «Гофермарт»
//	@BasePath		/

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/api/user/register", UserRegister)
	authorized := r.Group("/")
	authorized.Use(middleware.Auth())
	authorized.POST("/api/user/orders", PostUserOrders)
	authorized.GET("/api/user/orders", GetUserOrders)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
