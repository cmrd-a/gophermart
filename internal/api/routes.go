package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
)

// @Summary	Регистрация пользователя
// @Tags		user
// @Param		request	body	UserRegisterRequest	true	"данные пользователя для регистрации"
// @Accept		json
// @Produce	json
// @Success	200
// @Failure	400	{object}	httputil.HTTPError
// @Failure	409	{object}	httputil.HTTPError
// @Failure	500	{object}	httputil.HTTPError
// @Router		/api/user/register [post]
func UserRegister(c *gin.Context) {
	if true {
		c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
	} else {
		httputil.NewError(c, http.StatusBadRequest, nil)
	}
	return
}

// Получение списка загруженных номеров заказов
//
//	@Success	200	{object}	Orders
//
//	@Router		/api/user/orders [get]
func GetUserOrders(c *gin.Context) {
	// Implementation goes here
}
