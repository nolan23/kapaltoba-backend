package http

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/transaction"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseError struct {
	Message string `json:"message"`
}

type HttpTransactionHandler struct {
	TransactionUsecase transaction.Usecase
}

func NewTransactionHttpHandler(e *echo.Echo, ts transaction.Usecase) {
	handler := &HttpTransactionHandler{
		TransactionUsecase: ts,
	}
	e.GET("/transactions", handler.FetchTransaction)
	e.POST("/transaction", handler.Store)
	e.GET("/transaction/:id", handler.GetByID)
	e.PUT("/transaction/:id/pay", handler.Pay)
	e.PUT("/transaction/:id/cancel", handler.Cancel)
}

func (h *HttpTransactionHandler) FetchTransaction(c echo.Context) error {
	limit := c.QueryParam("limit")
	limitNum, _ := strconv.Atoi(limit)
	skip := c.QueryParam("skip")
	skipNum, _ := strconv.Atoi(skip)
	sort := c.QueryParam("sort")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listTransaction, nextSkip, err := h.TransactionUsecase.Fetch(ctx, limitNum, skipNum, sort)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Skip`, strconv.Itoa(nextSkip))
	return c.JSON(http.StatusOK, listTransaction)
}

func (h *HttpTransactionHandler) Store(c echo.Context) error {
	var transaction models.Transaction
	err := c.Bind(&transaction)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.TransactionUsecase.Store(ctx, &transaction)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, transaction)
}

func (h *HttpTransactionHandler) GetByID(c echo.Context) error {
	requestId := c.Param("id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	user, err := h.TransactionUsecase.GetByID(ctx, requestId)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *HttpTransactionHandler) Pay(c echo.Context) error {

	err := h.updateTransactionStatus(c, "lunas")
	if err != nil {
		log.Println("error update transaction to pay " + err.Error())
		return err
	}

	return nil
}

func (h *HttpTransactionHandler) Cancel(c echo.Context) error {

	err := h.updateTransactionStatus(c, "batal")
	if err != nil {
		log.Println("error update transaction to cancel " + err.Error())
		return err
	}

	return nil
}

func (h *HttpTransactionHandler) updateTransactionStatus(c echo.Context, status string) error {
	requestId := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		log.Println("error in handler " + err.Error())
		return err
	}
	ctx := c.Request().Context()
	if ctx == nil {
		log.Println("not found transaction in update transaction status")
		return nil
	}
	err = h.TransactionUsecase.Update(context.Background(), bson.M{"_id": oid}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		log.Println("error update handler " + err.Error())
		return err
	}
	return nil
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
