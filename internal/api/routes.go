package api

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/cmrd-a/gophermart/internal/repository"
	"github.com/cmrd-a/gophermart/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
)

var repo, _ = repository.NewRepository()
var Service = service.NewService(*repo)

// UserRegister регистрирует нового пользователя
//
//	@Summary	Регистрация пользователя
//	@Tags		auth
//	@Param		request	body	UserRegisterRequest	true	"данные пользователя для регистрации"
//	@Accept		json
//	@Produce	json
//	@Success	200	"пользователь успешно зарегистрирован и аутентифицирован"
//	@Failure	400	{object}	httputil.HTTPError	"неверный формат запроса"
//	@Failure	409	{object}	httputil.HTTPError	"логин уже занят"
//	@Failure	500	{object}	httputil.HTTPError	"внутренняя ошибка сервера"
//	@Header		200	string		Authorization		"токен авторизации"
//	@Router		/api/user/register [post]
func UserRegister(c *gin.Context) {
	r := UserRegisterRequest{}
	if err := c.BindJSON(&r); err != nil {
		c.String(http.StatusOK, err.Error())
	}
	userID, err := Service.AddUser(c, r.Login, r.Password)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	c.Header("Authorization", strconv.Itoa(int(userID)))
}

// PostUserOrders загружает номер заказа
//
//	@Summary	Загрузка номера заказа
//	@Success	200	"номер заказа уже был загружен этим пользователем"
//	@Success	202	"новый номер заказа принят в обработку"
//	@Failure	400	"неверный формат запроса"
//	@Failure	401	"пользователь не аутентифицирован"
//	@Failure	409	"номер заказа уже был загружен другим пользователем"
//	@Failure	422	"неверный формат номера заказа"
//	@Failure	500	"внутренняя ошибка сервера"
//	@Tags		orders
//	@Accept		plain
//	@Param		orderNumber		body	string	true	"номер заказа"
//	@Param		Authorization	header	string	true	"токен авторизации"
//	@Produce	json
//	@Router		/api/user/orders [post]
func PostUserOrders(c *gin.Context) {

	userID := c.GetInt64("userID")

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	orderNumber := string(bodyBytes)
	orderNumberInt, err := strconv.Atoi(orderNumber)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	if !service.Valid(orderNumberInt) {
		httputil.NewError(c, http.StatusUnprocessableEntity, errors.New("invalid order number"))
		return
	}
	existed := Service.GetOrder(c, orderNumber)
	if existed != nil {
		if existed.UserID == userID {
			c.Status(http.StatusOK)
			return
		} else {
			httputil.NewError(c, http.StatusConflict, errors.New("order already exists"))
			return
		}
	}
	err = Service.AddOrder(c, orderNumber, userID)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusAccepted)
}

// GetUserOrders возвращает список загруженных номеров заказов
//
//	@Summary	Получение списка загруженных номеров заказов
//	@Failure	401	{object}	httputil.HTTPError	"пользователь не авторизован"
//	@Failure	500	{object}	httputil.HTTPError	"внутренняя ошибка сервера"
//	@Tags		orders
//	@Success	200	{object}	Orders	"успешная обработка запроса"
//	@Success	204	{object}	Orders	"нет данных для ответа"
//	@Produce	json
//	@Param		Authorization	header	string	true	"токен авторизации"
//	@Router		/api/user/orders [get]
func GetUserOrders(c *gin.Context) {
	userID := c.GetInt64("userID")

	ro, err := Service.GetUserOrders(c, userID)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}
	if len(ro) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	jo := make(Orders, len(ro))
	for i, order := range ro {
		jo[i] = Order{Number: order.Number, Status: order.Status, Accural: order.Accural, UploadedAt: JSONTime(order.UploadedAt)}
	}
	c.JSON(http.StatusOK, &jo)
}
