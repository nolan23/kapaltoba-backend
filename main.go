package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/nolan23/kapaltoba-backend/models"
	_userHttpDeliver "github.com/nolan23/kapaltoba-backend/user/delivery/http"
	_userRepo "github.com/nolan23/kapaltoba-backend/user/repository"
	_userUsecase "github.com/nolan23/kapaltoba-backend/user/usecase"

	_transactionHttpDeliver "github.com/nolan23/kapaltoba-backend/transaction/delivery/http"
	_transactionRepo "github.com/nolan23/kapaltoba-backend/transaction/repository"
	_transactionUsecase "github.com/nolan23/kapaltoba-backend/transaction/usecase"

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
	userRepo := _userRepo.NewMongoDBUserRepository(con)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	userUsecase := _userUsecase.NewUserUsecase(userRepo, timeoutContext)
	_userHttpDeliver.NewUserHttpHandler(e, userUsecase)

	transactionRepo := _transactionRepo.NewMongoDBTransactionRepository(con)
	transactionUsecase := _transactionUsecase.NewTransactionUsecase(transactionRepo, timeoutContext)
	_transactionHttpDeliver.NewTransactionHttpHandler(e, transactionUsecase)
	var trn *models.User
	trn = &models.User{
		Name:         "roby",
		Email:        "test",
		PhoneNumber:  "phone",
		BirthDate:    time.Now(),
		Password:     "test123",
		ImageProfile: "imageprofile",
		TripHistory:  nil,
	}
	err = userRepo.Store(context.Background(), trn)
	fmt.Println("test")
	if err != nil {
		log.Println(err.Error())
		fmt.Println("error saving " + err.Error())
	}
	e.Start(viper.GetString("server.address"))
}
