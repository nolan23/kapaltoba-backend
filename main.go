package main

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/nolan23/kapaltoba-backend/models"
	_userHttpDeliver "github.com/nolan23/kapaltoba-backend/user/delivery/http"
	_userRepo "github.com/nolan23/kapaltoba-backend/user/repository"
	_userUsecase "github.com/nolan23/kapaltoba-backend/user/usecase"
	"github.com/spf13/viper"
	"github.com/zebresel-com/mongodm"
)

var serverMongo = viper.GetString(`database.host`) + ":" + viper.GetString(`database.port`)
var dbConfig = &mongodm.Config{
	DatabaseHosts:    []string{"127.0.0.1:27017"},
	DatabaseName:     viper.GetString(`database.name`),
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

	e.Start(viper.GetString("server.address"))

}
