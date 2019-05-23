package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/nolan23/kapaltoba-backend/boat"
	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseError struct {
	Message string `json:"message"`
}

type HttpBoatHandler struct {
	BoatUsecase boat.Usecase
}

func NewBoatHttpHandler(e *echo.Echo, bu boat.Usecase) {
	handler := &HttpBoatHandler{
		BoatUsecase: bu,
	}
	e.GET("/boats", handler.FetchBoat)
	e.POST("/boat", handler.Store)
	e.PUT("/boat:id", handler.Edit)
	e.GET("/boat/:id", handler.GetByID)
}

func (h *HttpBoatHandler) FetchBoat(c echo.Context) error {
	limit := c.QueryParam("limit")
	limitNum, _ := strconv.Atoi(limit)
	skip := c.QueryParam("skip")
	skipNum, _ := strconv.Atoi(skip)
	sort := c.QueryParam("sort")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listBoat, nextSkip, err := h.BoatUsecase.Fetch(ctx, limitNum, skipNum, sort)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Skip`, strconv.Itoa(nextSkip))
	return c.JSON(http.StatusOK, listBoat)
}

func (h *HttpBoatHandler) Store(c echo.Context) error {
	var boat models.Boat
	err := c.Bind(&boat)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.BoatUsecase.Store(ctx, &boat)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, boat)
}

func (h *HttpBoatHandler) GetByID(c echo.Context) error {
	requestId := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	user, err := h.BoatUsecase.GetByID(ctx, requestId)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *HttpBoatHandler) Edit(c echo.Context) error {
	requestId := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		log.Println("error in handler " + err.Error())
		return err
	}
	var boat models.Boat
	err = c.Bind(&boat)
	if err != nil {
		fmt.Println("you are error " + err.Error())
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	// bson.M{"boat": trip.Boat, "origin": trip.Origin, "destination": trip.Destination}
	boat.ID = oid
	err = h.BoatUsecase.Update(ctx, bson.M{"_id": oid}, &boat)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, boat)
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
