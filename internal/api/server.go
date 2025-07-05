package api

import (
	"log"

	_ "github.com/cmrd-a/gophermart/internal/api/docs"
	"github.com/cmrd-a/gophermart/internal/api/middleware"
	"github.com/cmrd-a/gophermart/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Gophermart API
//	@version		1.0
//	@description	Накопительная система лояльности «Гофермарт»
//	@BasePath		/

func SetupRouter(svc *service.Service) *gin.Engine {
	r := gin.Default()
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Printf("error on set trusted proxies:%e", err)
	}

	r.POST("/api/user/register", UserRegister(svc))
	r.POST("/api/user/login", UserLogin(svc))
	authorized := r.Group("/")
	authorized.Use(middleware.Auth())
	authorized.POST("/api/user/orders", PostUserOrders(svc))
	authorized.GET("/api/user/orders", GetUserOrders(svc))
	authorized.GET("/api/user/balance", GetUserBalance(svc))
	authorized.POST("/api/user/balance/withdraw", UserBalanceWithdraw(svc))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
