package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nolan23/kapaltoba-backend/credential"
	"github.com/nolan23/kapaltoba-backend/models"

	"github.com/labstack/echo/middleware"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	_userHttpDeliver "github.com/nolan23/kapaltoba-backend/user/delivery/http"
	_userRepo "github.com/nolan23/kapaltoba-backend/user/repository"
	_userUsecase "github.com/nolan23/kapaltoba-backend/user/usecase"

	_transactionHttpDeliver "github.com/nolan23/kapaltoba-backend/transaction/delivery/http"
	_transactionRepo "github.com/nolan23/kapaltoba-backend/transaction/repository"
	_transactionUsecase "github.com/nolan23/kapaltoba-backend/transaction/usecase"

	_tripHttpDeliver "github.com/nolan23/kapaltoba-backend/trip/delivery/http"
	_tripRepo "github.com/nolan23/kapaltoba-backend/trip/repository"
	_tripUsecase "github.com/nolan23/kapaltoba-backend/trip/usecase"

	_boatHttpDeliver "github.com/nolan23/kapaltoba-backend/boat/delivery/http"
	_boatRepo "github.com/nolan23/kapaltoba-backend/boat/repository"
	_boatUsecase "github.com/nolan23/kapaltoba-backend/boat/usecase"

	_credentialHttpDeliver "github.com/nolan23/kapaltoba-backend/credential/delivery/http"
	_credentialRepo "github.com/nolan23/kapaltoba-backend/credential/repository"
	_credentialUsecase "github.com/nolan23/kapaltoba-backend/credential/usecase"

	_captainHttpDeliver "github.com/nolan23/kapaltoba-backend/captain/delivery/http"
	_captainRepo "github.com/nolan23/kapaltoba-backend/captain/repository"
	_captainUsecase "github.com/nolan23/kapaltoba-backend/captain/usecase"

	"github.com/labstack/echo"

	"github.com/spf13/viper"
)

var serverMongo = viper.GetString(`database.host`) + ":" + viper.GetString(`database.port`)
var mongoURI = "mongodb+srv://roby:roby_is_the_best@cluster0-ld8yy.mongodb.net"
var mongoOld = "mongodb://roby:roby_is_the_best@cluster0-shard-00-00-ld8yy.mongodb.net:27017,cluster0-shard-00-01-ld8yy.mongodb.net:27017,cluster0-shard-00-02-ld8yy.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin&retryWrites=true"
var uri string

// var dbConfig = &mongodm.Config{
// 	DatabaseHosts:    []string{"mongodb://cluster0-shard-00-00-ld8yy.mongodb.net.:27017"},
// 	DatabaseName:     "kapaltoba",
// 	DatabaseUser:     viper.GetString(`database.user`),
// 	DatabasePassword: viper.GetString(`database.pass`),
// 	DatabaseSource:   "",
// }

// var dbConfig = &mongodm.Config{
// 	DialInfo:       dialInfo,
// 	DatabaseSource: "",
// }

type jwtCusctomClaimns struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// func login(c echo.Context) error {
// 	username := c.FormValue("username")
// 	password := c.FormValue("password")

// 	if username != "roby" || password != "roby123" {
// 		return echo.ErrUnauthorized
// 	}
// 	claims := &models.Claims{
// 		"Roby",
// 		"User",
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(time.Hour * 36).Unix(),
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	t, err := token.SignedString([]byte(viper.GetString("jwt.private")))
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, echo.Map{
// 		"token": t,
// 	})
// }

func restricted(c echo.Context) error {
	log.Println("masuk restricted")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.Claims)
	name := claims.Username
	return c.String(http.StatusOK, "welcome "+name)
}

func init() {
	var ok bool
	uri, ok = os.LookupEnv("MONGODB_URI")
	if !ok {
		uri = viper.GetString("database.uri")
	}
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
	uri = "mongodb://roby:roby123@localhost:27017/?authSource=admin"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	defer client.Disconnect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Connected to MongoDB!")
	database := client.Database(viper.GetString("database.namedev"))
	collection := database.Collection("test")
	_, err = collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	log.Println("Connected to MongoDB!")
	e := echo.New()
	r := e.Group("")
	config := middleware.JWTConfig{
		Claims:        &models.Claims{},
		SigningKey:    []byte(viper.GetString("jwt.private")),
		SigningMethod: "HS256",
	}
	r.Use(middleware.JWTWithConfig(config))
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// e.POST("/login", login)
	// r = e.Group("/restricted")
	// config := middleware.JWTConfig{
	// 	Claims:     &jwtCusctomClaimns{},
	// 	SigningKey: []byte("rahasia"),
	// }
	// r.Use(middleware.JWTWithConfig(config))
	// r.GET("/restricted", restricted)
	credential.Init()
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	userRepo := _userRepo.NewMongoUserRepository(database, "user")
	tripRepo := _tripRepo.NewMongoTripRepository(database, "trip")
	boatRepo := _boatRepo.NewMongoBoatRepository(database, "boat")
	credentialRepo := _credentialRepo.NewMongoCredentialRepository(database, "credential")
	transactionRepo := _transactionRepo.NewMongoTransactionRepository(database, "transaction")
	captainRepo := _captainRepo.NewMongoCaptainRepository(database, "captain")

	userUsecase := _userUsecase.NewUserUsecase(userRepo, tripRepo, transactionRepo, credentialRepo, timeoutContext)
	_userHttpDeliver.NewUserHttpHandler(e, userUsecase)

	transactionUsecase := _transactionUsecase.NewTransactionUsecase(transactionRepo, timeoutContext)
	_transactionHttpDeliver.NewTransactionHttpHandler(r, transactionUsecase)

	boatUsecase := _boatUsecase.NewBoatUsecase(boatRepo, captainRepo, timeoutContext)
	_boatHttpDeliver.NewBoatHttpHandler(e, boatUsecase)

	tripUsecase := _tripUsecase.NewTripUsecase(tripRepo, userRepo, boatRepo, captainRepo, timeoutContext)
	_tripHttpDeliver.NewTripHttpHandler(e, tripUsecase, boatUsecase)

	captainUsecase := _captainUsecase.NewCaptainUsecase(captainRepo, credentialRepo, tripRepo, timeoutContext)
	_captainHttpDeliver.NewCaptainHttpHandler(e, captainUsecase, tripUsecase, userUsecase, transactionUsecase)

	credentialUsecase := _credentialUsecase.NewCredentialUsecase(credentialRepo, timeoutContext)
	_credentialHttpDeliver.NewCredentialsHttpHandler(e, credentialUsecase, userUsecase, captainUsecase)

	port, ok := os.LookupEnv("PORT")

	if ok == false {
		port = "3000"
	}
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + port))
}
