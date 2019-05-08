package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/trip"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type HttpTripHandler struct {
	TripUsecase trip.Usecase
}

func NewTripHttpHandler(e *echo.Echo, ts trip.Usecase) {
	handler := &HttpTripHandler{
		TripUsecase: ts,
	}
	e.GET("/trips", handler.FetchTrip)
	e.POST("/trip", handler.Store)
	e.GET("/trip/:id/passengers", handler.GetPassenger)

}

func (h *HttpTripHandler) FetchTrip(c echo.Context) error {
	limit := c.QueryParam("limit")
	limitNum, _ := strconv.Atoi(limit)
	skip := c.QueryParam("skip")
	skipNum, _ := strconv.Atoi(skip)
	sort := c.QueryParam("sort")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listTrip, nextSkip, err := h.TripUsecase.Fetch(ctx, limitNum, skipNum, sort)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Skip`, strconv.Itoa(nextSkip))
	return c.JSON(http.StatusOK, listTrip)
}

func (h *HttpTripHandler) GetPassenger(c echo.Context) error {
	tripId := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	passengers, err := h.TripUsecase.GetPassenger(ctx, tripId)
	if err != nil {
		log.Println("error get passenger trip handler " + err.Error())
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, passengers)
}

func (h *HttpTripHandler) Store(c echo.Context) error {
	fmt.Println("you are here")
	var trip models.Trip
	err := c.Bind(&trip)
	if err != nil {
		fmt.Println("you are error " + err.Error())
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	err = h.TripUsecase.Store(ctx, &trip)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, trip)
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
