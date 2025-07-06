package api

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/cmrd-a/gophermart/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
	"github.com/theplant/luhn"
)

// UserRegister регистрирует нового пользователя
//
//	@Summary	Регистрация пользователя
//	@Tags		auth
//	@Param		request	body	LoginPasswordRequest	true	"данные пользователя для регистрации"
//	@Accept		json
//	@Produce	json
//	@Success	200	"пользователь успешно зарегистрирован и аутентифицирован"
//	@Failure	400	{object}	httputil.HTTPError	"неверный формат запроса"
//	@Failure	409	{object}	httputil.HTTPError	"логин уже занят"
//	@Failure	500	{object}	httputil.HTTPError	"внутренняя ошибка сервера"
//	@Header		200	string		Authorization		"токен авторизации"
//	@Router		/api/user/register [post]
func UserRegister(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		r := LoginPasswordRequest{}
		if err := c.BindJSON(&r); err != nil {
			c.String(http.StatusOK, err.Error())
		}
		userID, err := svc.AddUser(c, r.Login, r.Password)
		if err != nil {
			httputil.NewError(c, http.StatusBadRequest, err)
			return
		}
		token, err := BuildJWTString(userID)
		if err != nil {
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}
		c.Header("Authorization", token)
	}
}

// UserLogin регистрирует нового пользователя
//
//	@Summary	Аутентификация пользователя
//	@Tags		auth
//	@Param		request	body	LoginPasswordRequest	true	"данные пользователя для регистрации"
//	@Accept		json
//	@Produce	json
//	@Success	200	"пользователь успешно аутентифицирован"
//	@Failure	400	{object}	httputil.HTTPError	"неверный формат запроса"
//	@Failure	401	{object}	httputil.HTTPError	"неверная пара логин/пароль"
//	@Failure	500	{object}	httputil.HTTPError	"внутренняя ошибка сервера"
//	@Header		200	string		Authorization		"токен авторизации"
//	@Router		/api/user/login [post]
func UserLogin(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		r := LoginPasswordRequest{}
		if err := c.BindJSON(&r); err != nil {
			c.String(http.StatusOK, err.Error())
		}
		userID, err := svc.CheckLoginPassword(c, r.Login, r.Password)
		if err != nil {
			httputil.NewError(c, http.StatusBadRequest, err)
			return
		}
		if userID > 0 {
			token, err := BuildJWTString(userID)
			if err != nil {
				httputil.NewError(c, http.StatusInternalServerError, err)
				return
			}
			c.Header("Authorization", token)
			return
		}
		httputil.NewError(c, http.StatusUnauthorized, errors.New("invalid login/password pair"))
	}
}

// PostUserOrders загружает номер заказа
//
//	@Summary	Загрузка номера заказа
//	@Tags		orders
//	@Success	200	"номер заказа уже был загружен этим пользователем"
//	@Success	202	"новый номер заказа принят в обработку"
//	@Failure	400	"неверный формат запроса"
//	@Failure	401	"пользователь не аутентифицирован"
//	@Failure	409	"номер заказа уже был загружен другим пользователем"
//	@Failure	422	"неверный формат номера заказа"
//	@Failure	500	"внутренняя ошибка сервера"
//	@Accept		plain
//	@Param		orderNumber		body	string	true	"номер заказа"
//	@Param		Authorization	header	string	true	"токен авторизации"
//	@Produce	json
//	@Router		/api/user/orders [post]
func PostUserOrders(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
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
		if !luhn.Valid(orderNumberInt) {
			httputil.NewError(c, http.StatusUnprocessableEntity, errors.New("invalid order number"))
			return
		}
		existed := svc.GetOrder(c, orderNumber)
		if existed != nil {
			if existed.UserID == userID {
				c.Status(http.StatusOK)
				return
			} else {
				httputil.NewError(c, http.StatusConflict, errors.New("order already exists"))
				return
			}
		}
		err = svc.AddOrder(c, orderNumber, userID)
		if err != nil {
			httputil.NewError(c, http.StatusBadRequest, err)
			return
		}
		svc.Publish(orderNumber)
		c.Status(http.StatusAccepted)
	}
}

// GetUserOrders возвращает список загруженных номеров заказов
//
//	@Summary	Получение списка загруженных номеров заказов
//	@Tags		orders
//	@Failure	401	{object}	httputil.HTTPError		"пользователь не авторизован"
//	@Failure	500	{object}	httputil.HTTPError		"внутренняя ошибка сервера"
//	@Success	200	{object}	GetUserOrdersResponse	"успешная обработка запроса"
//	@Success	204	"нет данных для ответа"
//	@Produce	json
//	@Param		Authorization	header	string	true	"токен авторизации"
//	@Router		/api/user/orders [get]
func GetUserOrders(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		userID := c.GetInt64("userID")
		ro, err := svc.GetUserOrders(c, userID)
		if err != nil {
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}
		if len(ro) == 0 {
			c.Status(http.StatusNoContent)
			return
		}
		res := make(GetUserOrdersResponse, len(ro))
		for i, order := range ro {
			res[i] = Order{
				Number:     order.Number,
				Status:     order.Status,
				Accrual:    order.Accrual.InexactFloat64(),
				UploadedAt: JSONTime(order.UploadedAt),
			}
		}
		c.JSON(http.StatusOK, &res)
	}
}

// GetUserBalance возвращает текущий баланс пользователя
//
//	@Summary	Получение текущего баланса пользователя
//	@Tags		balance
//	@Failure	401	{object}	httputil.HTTPError	"пользователь не авторизован"
//	@Failure	500	{object}	httputil.HTTPError	"внутренняя ошибка сервера"
//	@Success	200	{object}	BalanceResponse		"успешная обработка запроса"
//	@Produce	json
//	@Param		Authorization	header	string	true	"токен авторизации"
//	@Router		/api/user/balance [get]
func GetUserBalance(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		userID := c.GetInt64("userID")
		balance, err := svc.GetUserBalance(c, userID)
		if err != nil {
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}
		jb := BalanceResponse{Current: balance.Current.InexactFloat64(), Withdrawn: balance.Withdraw.InexactFloat64()}
		c.JSON(http.StatusOK, &jb)
	}
}

// UserBalanceWithdraw списывает средства
//
//	@Summary	Запрос на списание средств
//	@Tags		balance
//	@Failure	400	{object}	httputil.HTTPError	"неверный формат запроса"
//	@Failure	401	{object}	httputil.HTTPError	"пользователь не аутентифицирован"
//	@Failure	409	{object}	httputil.HTTPError	"номер заказа уже был загружен другим пользователем"
//	@Failure	422	{object}	httputil.HTTPError	"неверный формат номера заказа"
//	@Failure	402	{object}	httputil.HTTPError	"на счету недостаточно средств"
//	@Failure	402	{object}	httputil.HTTPError	"неверный номер заказа"
//	@Failure	500	{object}	httputil.HTTPError	"внутренняя ошибка сервера"
//	@Success	200	"успешная обработка запроса"
//	@Accept		json
//	@Produce	json
//	@Param		request			body	UserBalanceWithdrawRequest	true	"номер счёта и сумма баллов к списанию в счёт оплаты"
//	@Param		Authorization	header	string						true	"токен авторизации"
//	@Router		/api/user/balance/withdraw [POST]
func UserBalanceWithdraw(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		userID := c.GetInt64("userID")
		r := UserBalanceWithdrawRequest{}
		if err := c.BindJSON(&r); err != nil {
			c.String(http.StatusOK, err.Error())
		}
		orderNumberInt, err := strconv.Atoi(r.Order)
		if err != nil {
			httputil.NewError(c, http.StatusBadRequest, err)
			return
		}
		if !luhn.Valid(orderNumberInt) {
			httputil.NewError(c, http.StatusUnprocessableEntity, errors.New("invalid order number"))
			return
		}
		err = svc.WithdrawUserBalance(c, r.Order, userID, decimal.NewFromFloat(r.Sum))
		if err != nil {
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}
	}
}

// GetUserWithdrawals возвращает список выводов средств
//
//	@Summary	Получение информации о выводе средств
//	@Tags		balance
//	@Success	200	{object}	GetUserWithdrawalsResponse	"успешная обработка запроса"
//	@Success	204	"нет ни одного списания"
//	@Failure	401	{object}	httputil.HTTPError	"пользователь не авторизован"
//	@Failure	500	{object}	httputil.HTTPError	"внутренняя ошибка сервера"
//	@Produce	json
//	@Param		Authorization	header	string	true	"токен авторизации"
//	@Router		/api/user/withdrawals [get]
func GetUserWithdrawals(svc *service.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		userID := c.GetInt64("userID")
		withdrawals, err := svc.GetUserWithdrawals(c, userID)
		if err != nil {
			httputil.NewError(c, http.StatusInternalServerError, err)
			return
		}
		if len(withdrawals) == 0 {
			c.Status(http.StatusNoContent)
			return
		}
		res := make(GetUserWithdrawalsResponse, len(withdrawals))
		for i, withdraw := range withdrawals {
			res[i] = Withdraw{
				Order:       withdraw.OrderNumber,
				Sum:         withdraw.Withdraw.InexactFloat64(),
				ProcessedAt: JSONTime(withdraw.ProcessedAt),
			}
		}
		c.JSON(http.StatusOK, &res)
	}
}
