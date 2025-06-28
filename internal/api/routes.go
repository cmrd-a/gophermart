package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
)

// UserRegister регистрирует нового пользователя
//
//	@Summary	Регистрация пользователя
//	@Tags		auth
//	@Param		request	body	UserRegisterRequest	true	"данные пользователя для регистрации"
//	@Accept		json
//	@Produce	json
//	@Success	200
//	@Failure	400	{object}	httputil.HTTPError
//	@Failure	409	{object}	httputil.HTTPError
//	@Failure	500	{object}	httputil.HTTPError
//	@Router		/api/user/register [post]
func UserRegister(c *gin.Context) {
	r := UserRegisterRequest{}
	if err := c.BindJSON(&r); err != nil {
		c.String(http.StatusOK, err.Error())
	}
	err := fmt.Errorf("adasd")
	if err == nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	// c.SetCookie("gin_cookie", "test1", 3600, "/", "localhost", false, true)
	c.Header("Authorization", r.Login)

}

// PostUserOrders загружает номер заказа
//
//	@Summary	Загрузка номера заказа
//	@Tags		orders
//	@Accept		text/plain
//	@Produce	text/plain
//	@Router		/api/user/orders [post]
func PostUserOrders(c *gin.Context) {
	// 200 — номер заказа уже был загружен этим пользователем;
	// 202 — новый номер заказа принят в обработку;
	// 400 — неверный формат запроса;
	// 401 — пользователь не аутентифицирован;
	// 409 — номер заказа уже был загружен другим пользователем;
	// 422 — неверный формат номера заказа;
	// 500 — внутренняя ошибка сервера.
	user := c.GetString("userLogin")

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	s := string(bodyBytes)
	on, err := strconv.Atoi(s)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	o, ok := db[user]
	if ok && o == int64(on) {
		c.Status(http.StatusOK)
		return
	}
	db[user] = int64(on)
	c.Status(http.StatusAccepted)
}

// GetUserOrders возвращает список загруженных номеров заказов
//
//	@Summary	Получение списка загруженных номеров заказов
//	@Tags		orders
//	@Success	200	{object}	Orders
//	@Produce	application/json
//	@Router		/api/user/orders [get]
func GetUserOrders(c *gin.Context) {
	// 200 Orders
	// 204 — нет данных для ответа.
	// 401 — пользователь не авторизован.
	// 500 — внутренняя ошибка сервера.

	user := c.GetString("userLogin")
	orders := make(Orders, 0)
	o, ok := db[user]
	if ok {
		order := Order{Number: o}
		orders = append(orders, order)
	}
	c.JSON(http.StatusOK, &orders)
}
