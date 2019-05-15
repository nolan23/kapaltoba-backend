package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nolan23/kapaltoba-backend/credential"

	"github.com/labstack/echo"
	"github.com/nolan23/kapaltoba-backend/models"
)

type ResponseError struct {
	Message string `json:"message"`
}

type HttpCredentialHandler struct {
	CredentialUsecase credential.Usecase
}

func NewCredentialsHttpHandler(e *echo.Echo, credentialUsecase credential.Usecase) {
	handler := &HttpCredentialHandler{
		CredentialUsecase: credentialUsecase,
	}
	e.POST("/signin", handler.SignIn)
	e.POST("/signup", handler.SignUp)
	e.POST("/signout", handler.SignOut)
}

func (h *HttpCredentialHandler) SignIn(c echo.Context) error {
	log.Println("masuk")
	var cred models.Credential
	err := c.Bind(&cred)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": err.Error()})
	}

	if cred.Username == "" {
		log.Println("Username is required")
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Username is required"})
	}
	if cred.Password == "" {
		log.Println("Password is required")
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Password is required"})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	credAux, err := h.CredentialUsecase.GetByUsername(ctx, cred.Username)
	if err != nil {
		log.Println("failed to get credential by username " + err.Error())
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}
	if credAux.Username == "" {
		log.Println(fmt.Sprintf("User %s not registered.", cred.Username))
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "User not registered"})
	}

	err = credential.CompareHashedPasswords(credAux.Password, cred.Password)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	tokenString, err := credential.GenerateJWT(credAux.Username, credAux.Role)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	// authCookie := http.Cookie{
	// 	Name:     "AuthToken",
	// 	Value:    tokenString,
	// 	HttpOnly: true,
	// }
	// c.SetCookie(&authCookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})
}

func (h *HttpCredentialHandler) SignUp(c echo.Context) error {
	var cred models.Credential
	err := c.Bind(&cred)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if cred.Username == "" {
		log.Println("Username is required")
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Username is required"})
	}
	if cred.Password == "" {
		log.Println("Password is required")
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Password is required"})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	var credAux *models.Credential
	credAux, err = h.CredentialUsecase.GetByUsername(ctx, cred.Username)

	// if err != nil {
	// 	log.Println(err)
	// 	return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	// }
	if credAux != nil {
		if cred.Username == credAux.Username {
			log.Println(fmt.Sprintf("Username %s already exists.", cred.Username))
			return c.JSON(http.StatusBadRequest, ResponseError{Message: "Username already exists"})
		}
	}
	hashedPassword, err := credential.GenerateHashedPassword(cred.Password)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusForbidden, ResponseError{Message: "Password generation failed"})
	}

	cred.Password = string(hashedPassword)
	err = h.CredentialUsecase.Store(ctx, &cred)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, ResponseError{Message: "User created successfully"})
}

func NullifyTokenCookies(c echo.Context) (string, error) {

	// If present, revoke the cookie.
	AuthCookie, err := c.Cookie("AuthToken")

	if err != nil {
		return "", err
	}

	// Remove the user's ability to make requests.
	jti, _ := credential.GrabJTI(AuthCookie.Value)

	credential.RevokeJWT(jti)

	// Set new authCookie without any token string.
	authCookie := http.Cookie{
		Name:     "AuthToken",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}
	c.SetCookie(&authCookie)

	return jti, nil
}

func (h *HttpCredentialHandler) SignOut(c echo.Context) error {

	// jti, err := NullifyTokenCookies(c)
	// if err == http.ErrNoCookie {
	// 	log.Println("No User logged in.")
	// 	return c.JSON(http.StatusBadRequest, ResponseError{Message: "No User logged in."})
	// }

	// if err != nil {
	// 	log.Println(err)
	// 	return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	// }

	return c.JSON(http.StatusOK, ResponseError{Message: "User  has been logged out"})
}
