package api

import (
	"log"

	_ "github.com/cmrd-a/gophermart/internal/api/docs"
	"github.com/cmrd-a/gophermart/internal/service"
	"github.com/gin-contrib/gzip"
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
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	err := r.SetTrustedProxies(nil)
	if err != nil {
		log.Printf("error on set trusted proxies:%e", err)
	}

	r.POST("/api/user/register", UserRegister(svc))
	r.POST("/api/user/login", UserLogin(svc))
	authorized := r.Group("/api/user")
	authorized.Use(Auth())
	authorized.POST("/orders", PostUserOrders(svc))
	authorized.GET("/orders", GetUserOrders(svc))
	authorized.GET("/balance", GetUserBalance(svc))
	authorized.POST("/balance/withdraw", UserBalanceWithdraw(svc))
	authorized.GET("/withdrawals", GetUserWithdrawals(svc))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
