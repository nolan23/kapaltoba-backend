package credential

import (
	"errors"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

const (
	authcJWTValidTime = 120 // authc JWT remains valid for 10 minutes.
)

var (
	privateKey string
	publicKey  string
)
var jwtWhiteList map[string]string

func Init() {
	log.Println("set private key ")
	privateKey = viper.GetString("jwt.private")
	log.Println("set private key " + privateKey)
	publicKey = viper.GetString("jwt.public")
	jwtWhiteList = make(map[string]string)
}

func GenerateHashedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CompareHashedPasswords(password1 string, password2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
}

func IsJTIInWhitelist(jti string) (string, bool) {

	jtw, exists := jwtWhiteList[jti]

	return jtw, exists
}

func GenerateJWT(id string, username string, name string, role string) (string, error) {

	// If there is a jwt in the whitelist for the user, invalidates it.
	_, ok := IsJTIInWhitelist(username)

	if ok {
		RevokeJWT(username)
	}

	claims := &models.Claims{
		ID:       id,
		Name:     name,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * authcJWTValidTime).Unix(),
		},
	}

	signer := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := signer.SignedString([]byte(privateKey))
	log.Println("generate with private key " + privateKey)
	if err != nil {
		return "", err
	}

	jwtWhiteList[username] = tokenString

	return tokenString, nil
}

func KeyFunc(token *jwt.Token) (interface{}, error) {

	return publicKey, nil
}

func ValidateJWT(tokenString string) error {

	_, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, KeyFunc)

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			// A few jwt validation errors.
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				return errors.New("Token validity expired.")
			case jwt.ValidationErrorSignatureInvalid:
				return errors.New("Token signature validation failed.")
			default:
				return errors.New("Token validation failed.")
			}
		}
	}

	return nil
}

func GrabJTI(tokenString string) (string, error) {

	token, _ := jwt.ParseWithClaims(tokenString, &models.Claims{}, KeyFunc)

	tokenClaims, ok := token.Claims.(*models.Claims)

	if !ok {
		return "", errors.New("Error when grabing jti from claims.")
	}

	// Grab jti from claims.
	jti := tokenClaims.StandardClaims.Subject

	return jti, nil
}

func RevokeJWT(jti string) {

	_, ok := IsJTIInWhitelist(jti)

	if ok {
		// Remove the JWT from the jwtWhiteList.
		delete(jwtWhiteList, jti)
	}
}
