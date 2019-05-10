package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/nolan23/kapaltoba-backend/models"
	_userHttpDeliver "github.com/nolan23/kapaltoba-backend/user/delivery/http"
	_userRepo "github.com/nolan23/kapaltoba-backend/user/repository"
	_userUsecase "github.com/nolan23/kapaltoba-backend/user/usecase"

	_transactionHttpDeliver "github.com/nolan23/kapaltoba-backend/transaction/delivery/http"
	_transactionRepo "github.com/nolan23/kapaltoba-backend/transaction/repository"
	_transactionUsecase "github.com/nolan23/kapaltoba-backend/transaction/usecase"

	_tripHttpDeliver "github.com/nolan23/kapaltoba-backend/trip/delivery/http"
	_tripRepo "github.com/nolan23/kapaltoba-backend/trip/repository"
	_tripUsecase "github.com/nolan23/kapaltoba-backend/trip/usecase"

	"github.com/spf13/viper"
	"github.com/zebresel-com/mongodm"
)

var serverMongo = viper.GetString(`database.host`) + ":" + viper.GetString(`database.port`)
var dbConfig = &mongodm.Config{
	DatabaseHosts:    []string{"127.0.0.1:27017"},
	DatabaseName:     "kapaltoba",
	DatabaseUser:     viper.GetString(`database.user`),
	DatabasePassword: viper.GetString(`database.pass`),
	DatabaseSource:   "",
}

type jwtCusctomClaimns struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "roby" || password != "roby123" {
		return echo.ErrUnauthorized
	}
	claims := &jwtCusctomClaimns{
		"Roby",
		"User",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 36).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("rahasia"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCusctomClaimns)
	name := claims.Name
	return c.String(http.StatusOK, "welcome "+name)
}

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}

}

func main() {
	var con, err = mongodm.Connect(dbConfig)
	con.Register(&models.User{}, "user")
	con.Register(&models.Transaction{}, "transaction")
	con.Register(&models.Trip{}, "trip")
	con.Register(&models.Boat{}, "boat")
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/login", login)
	r := e.Group("/restricted")
	config := middleware.JWTConfig{
		Claims:     &jwtCusctomClaimns{},
		SigningKey: []byte("rahasia"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("", restricted)
	userRepo := _userRepo.NewMongoDBUserRepository(con)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	userUsecase := _userUsecase.NewUserUsecase(userRepo, timeoutContext)
	_userHttpDeliver.NewUserHttpHandler(e, userUsecase)

	transactionRepo := _transactionRepo.NewMongoDBTransactionRepository(con)
	transactionUsecase := _transactionUsecase.NewTransactionUsecase(transactionRepo, timeoutContext)
	_transactionHttpDeliver.NewTransactionHttpHandler(e, transactionUsecase)

	tripRepo := _tripRepo.NewMongoDBTripRepository(con)
	tripUsecase := _tripUsecase.NewTripUsecase(tripRepo, userUsecase, timeoutContext)
	_tripHttpDeliver.NewTripHttpHandler(e, tripUsecase)
	// var trn *models.User
	// trn = &models.User{
	// 	Name:         "roby",
	// 	Email:        "test",
	// 	PhoneNumber:  "phone",
	// 	BirthDate:    time.Now(),
	// 	Password:     "test123",
	// 	ImageProfile: "imageprofile",
	// 	TripHistory:  nil,
	// }
	// err = userRepo.Store(context.Background(), trn)
	// fmt.Println("test")
	// if err != nil {
	// 	log.Println(err.Error())
	// 	fmt.Println("error saving " + err.Error())
	// }

	e.Start(viper.GetString("server.address"))
}
