package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

func init() {
	var ok bool
	uri, ok = os.LookupEnv("MONGODB_URI")
	if !ok {
		uri = viper.GetString("database.uri")
	}

}

func main() {
	// dbConfig := &mongodm.Config{
	// 	DialInfo:       dialInfo,
	// 	DatabaseSource: "",
	// }
	// uri = "mongodb://roby:roby123@localhost:27017/?authSource=admin"
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
	database := client.Database(viper.GetString("database.nameprod"))
	collection := database.Collection("test")
	_, err = collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	// 5cd6604db38f65c477040246
	// oid, erId := primitive.ObjectIDFromHex("5cd6604db38f65c477040246")
	// if erId != nil {
	// 	log.Println(erId.Error())
	// }
	// filter := bson.M{"_id": oid}
	// _, err = collection.InsertOne(ctx, bson.M{"name": "test", "value": "test"})
	// var result struct {
	// 	Name  string
	// 	Value string
	// }
	// err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	// log.Println(result.Name + " " + result.Value)
	log.Println("Connected to MongoDB!")

	// var con, err = mongodm.Connect(dbConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// 	os.Exit(1)
	// }

	// log.Println("Connected to MongoDB! 2")
	// con.Register(&models.User{}, "user")
	// con.Register(&models.Transaction{}, "transaction")
	// con.Register(&models.Trip{}, "trip")
	// con.Register(&models.Boat{}, "boat")

	// defer con.Close()
	e := echo.New()
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// e.POST("/login", login)
	// r := e.Group("/restricted")
	// config := middleware.JWTConfig{
	// 	Claims:     &jwtCusctomClaimns{},
	// 	SigningKey: []byte("rahasia"),
	// }
	// r.Use(middleware.JWTWithConfig(config))
	// r.GET("", restricted)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	userRepo := _userRepo.NewMongoUserRepository(database, "user")
	tripRepo := _tripRepo.NewMongoTripRepository(database, "trip")
	boatRepo := _boatRepo.NewMongoBoatRepository(database, "boat")
	transactionRepo := _transactionRepo.NewMongoTransactionRepository(database, "transaction")

	userUsecase := _userUsecase.NewUserUsecase(userRepo, tripRepo, transactionRepo, timeoutContext)
	_userHttpDeliver.NewUserHttpHandler(e, userUsecase)

	transactionUsecase := _transactionUsecase.NewTransactionUsecase(transactionRepo, timeoutContext)
	_transactionHttpDeliver.NewTransactionHttpHandler(e, transactionUsecase)

	tripUsecase := _tripUsecase.NewTripUsecase(tripRepo, userRepo, timeoutContext)
	_tripHttpDeliver.NewTripHttpHandler(e, tripUsecase)

	boatUsecase := _boatUsecase.NewBoatUsecase(boatRepo, timeoutContext)
	_boatHttpDeliver.NewBoatHttpHandler(e, boatUsecase)
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

	// e.Start(viper.GetString("server.address"))
	port, ok := os.LookupEnv("PORT")

	if ok == false {
		port = "3000"
	}
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + port))
}
