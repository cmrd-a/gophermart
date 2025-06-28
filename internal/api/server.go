package api

import (
	"net/http"

	_ "github.com/cmrd-a/gophermart/internal/api/docs"
	"github.com/cmrd-a/gophermart/internal/api/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Gophermart API
//	@version		1.0
//	@description	Накопительная система лояльности «Гофермарт»
//	@host			localhost:8080
//	@BasePath		/

var db = make(map[string]int64)

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.POST("/api/user/register", UserRegister)

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	authorized := r.Group("/")
	authorized.Use(middleware.Auth())
	authorized.POST("/api/user/orders", PostUserOrders)
	authorized.GET("/api/user/orders", GetUserOrders)

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value int64 `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
