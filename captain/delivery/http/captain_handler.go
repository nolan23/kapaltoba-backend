package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/nolan23/kapaltoba-backend/transaction"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nolan23/kapaltoba-backend/trip"
	"github.com/nolan23/kapaltoba-backend/user"

	"github.com/labstack/echo"
	"github.com/nolan23/kapaltoba-backend/captain"
	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ReturnData struct {
	Trip      *models.Trip `json:"trip"`
	Passenger []*Passenger `json:"passengers"`
}

type Passenger struct {
	PassengerID   string `json:"passengerID"`
	PassengerName string `json:"passengerName"`
	TransactionID string `json:"transactionID"`
	Status        string `json:"status"`
}

type HttpCaptainHandler struct {
	CaptainUsecase     captain.Usecase
	TripUsecase        trip.Usecase
	UserUsecase        user.Usecase
	TransactionUsecase transaction.Usecase
}

func NewCaptainHttpHandler(e *echo.Echo, bu captain.Usecase, tu trip.Usecase, uu user.Usecase, trs transaction.Usecase) {
	handler := &HttpCaptainHandler{
		CaptainUsecase:     bu,
		TripUsecase:        tu,
		UserUsecase:        uu,
		TransactionUsecase: trs,
	}
	e.GET("/captains", handler.FetchCaptain)
	e.POST("/captain", handler.Store)
	e.PUT("/captain/:id", handler.Edit)
	e.GET("/captain/:id", handler.GetByID)
	e.GET("/captain/u/:username", handler.GetByUsername)
	e.GET("/captain/:id/trips", handler.GetTrips)
	e.GET("/captain/trip/:idTrip", handler.GetTripDetail)
}

func (h *HttpCaptainHandler) FetchCaptain(c echo.Context) error {
	limit := c.QueryParam("limit")
	limitNum, _ := strconv.Atoi(limit)
	skip := c.QueryParam("skip")
	skipNum, _ := strconv.Atoi(skip)
	sort := c.QueryParam("sort")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listCaptain, nextSkip, err := h.CaptainUsecase.Fetch(ctx, limitNum, skipNum, sort)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Skip`, strconv.Itoa(nextSkip))
	return c.JSON(http.StatusOK, listCaptain)
}

func (h *HttpCaptainHandler) Store(c echo.Context) error {
	var captain models.Captain
	err := c.Bind(&captain)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.CaptainUsecase.Store(ctx, &captain)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, captain)
}

func (h *HttpCaptainHandler) Edit(c echo.Context) error {
	requestId := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		log.Println("error in handler " + err.Error())
		return err
	}
	var captain models.Captain
	err = c.Bind(&captain)
	if err != nil {
		fmt.Println("you are error " + err.Error())
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	captain.ID = oid
	err = h.CaptainUsecase.Update(ctx, bson.M{"_id": oid}, &captain)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, captain)
}

func (h *HttpCaptainHandler) GetByID(c echo.Context) error {
	requestId := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	user, err := h.CaptainUsecase.GetByID(ctx, requestId)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *HttpCaptainHandler) GetByUsername(c echo.Context) error {
	requestName := c.Param("username")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	captain, err := h.CaptainUsecase.GetByUsername(ctx, requestName)

	if captain == nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: "Captain not found"})
	}

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, captain)
}

func (h *HttpCaptainHandler) GetTrips(c echo.Context) error {
	requestId := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	trips, err := h.CaptainUsecase.GetTrips(ctx, requestId)
	if trips == nil {
		log.Println("null trips")
		return c.JSON(http.StatusNotFound, ResponseError{Message: "Trips not found"})
	}
	if err != nil {
		log.Println("error in get trips in captain handler " + err.Error())
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, trips)

}

func (h *HttpCaptainHandler) GetTripDetail(c echo.Context) error {
	idTrip := c.Param("idTrip")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	trip, err := h.TripUsecase.GetByID(ctx, idTrip)
	if trip == nil {
		log.Println("trip is empty ")
		return c.JSON(http.StatusNotFound, ResponseError{Message: "Trip Not Found"})
	}
	if err != nil {
		log.Println("error get trip " + err.Error())
		return c.JSON(http.StatusNotFound, ResponseError{Message: "Not Found"})
	}

	var res *ReturnData
	var passengers []*Passenger

	for _, user := range trip.Passengers {
		passenger, err := h.UserUsecase.GetByID(ctx, user)
		if err != nil {
			log.Println("error get user in trip usecase " + err.Error())
			continue
		}
		trans, er := h.TransactionUsecase.FindBy(ctx, passenger.ID.Hex(), trip.ID.Hex())
		if er != nil {
			log.Println("error get transaction " + er.Error())
			continue
		}
		var pas *Passenger
		pas = &Passenger{}
		pas.PassengerID = passenger.ID.Hex()
		pas.PassengerName = passenger.Name
		pas.TransactionID = trans.ID.Hex()
		pas.Status = trans.Status
		passengers = append(passengers, pas)
	}
	res = &ReturnData{}
	res.Trip = trip
	res.Passenger = passengers
	return c.JSON(http.StatusFound, res)
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
