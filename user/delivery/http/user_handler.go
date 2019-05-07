package http

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/user"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

// HttpUserHandler  represent the httphandler for user
type HttpUserHandler struct {
	UserUsecase user.Usecase
}

func NewUserHttpHandler(e *echo.Echo, us user.Usecase) {
	handler := &HttpUserHandler{
		UserUsecase: us,
	}
	e.GET("/users", handler.FetchUser)
	e.POST("/user", handler.Store)
	e.GET("user/:id", handler.GetByID)
	e.GET("user/:id/trips", handler.GetTrips)
}

func (h *HttpUserHandler) FetchUser(c echo.Context) error {
	limit := c.QueryParam("limit")
	limitNum, _ := strconv.Atoi(limit)
	skip := c.QueryParam("skip")
	skipNum, _ := strconv.Atoi(skip)
	sort := c.QueryParam("sort")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listUser, nextSkip, err := h.UserUsecase.Fetch(ctx, limitNum, skipNum, sort)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Skip`, strconv.Itoa(nextSkip))
	return c.JSON(http.StatusOK, listUser)
}

func (h *HttpUserHandler) Store(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.UserUsecase.Store(ctx, &user)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *HttpUserHandler) GetByID(c echo.Context) error {
	requestId := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	user, err := h.UserUsecase.GetByID(ctx, requestId)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *HttpUserHandler) GetTrips(c echo.Context) error {
	requestId := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	trips, err := h.UserUsecase.GetUserTrips(ctx, requestId)
	if err != nil {
		log.Println("error in get trips in user handler " + err.Error())
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, trips)
}

func getStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
