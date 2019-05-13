package repository

import (
	"context"
	"log"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/trip"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type mongoDBTripRepository struct {
	Conn *mongodm.Connection
}

func NewMongoDBTripRepository(Conn *mongodm.Connection) trip.Repository {
	return nil
}

func (ts *mongoDBTripRepository) fetch(ctx context.Context, query interface{}, limit int, skip int, sort string) (res []*models.Trip, nextSkip int, err error) {
	Trip := ts.Conn.Model("Trip")
	trip := []*models.Trip{}
	err = Trip.Find(query).Skip(skip).Limit(limit).Sort(sort).Exec(&trip)
	if err != nil {
		log.Println("error in find trip " + err.Error())
		return nil, skip + limit, err
	}
	return trip, skip + limit, nil
}
func (ts *mongoDBTripRepository) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Trip, nextSkip int, err error) {
	query := bson.M{"deleted": false}
	if sort == "" {
		sort = "_modifiedAt"
	}
	result, nextSkip, err := ts.fetch(ctx, query, limit, skip, sort)
	if err != nil {
		log.Println("error in fetch trip " + err.Error())
		return nil, nextSkip, err
	}
	return result, nextSkip, nil
}
func (ts *mongoDBTripRepository) GetByID(ctx context.Context, id string) (*models.Trip, error) {
	Trip := ts.Conn.Model("Trip")
	trip := &models.Trip{}
	err := Trip.FindId(bson.ObjectIdHex(id)).Populate("Passenger").Exec(trip)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		log.Println("error in get id trip " + err.Error())
		return nil, err
	} else if err != nil {
		log.Println("another error in fetch trip " + err.Error())
		return nil, err
	}

	// log.Println("trip id ", trip.Id.Hex(), " punya pass ", len(trip.Passenger.([]*models.User)))

	return trip, nil
}
func (ts *mongoDBTripRepository) Update(ctx context.Context, selector interface{}, update interface{}) error {
	Trip := ts.Conn.Model("Trip")
	err := Trip.Update(selector, update)
	if err != nil {
		log.Println("error update Trip " + err.Error())
		return err
	}
	return nil
}
func (ts *mongoDBTripRepository) Store(ctx context.Context, trip *models.Trip) error {
	// Trip := ts.Conn.Model("Trip")
	// Trip.New(trip)
	// err := trip.Save()

	// if err != nil {
	// 	log.Println("error in stre transaction " + err.Error())
	// 	return err
	// }
	// err = trip.Populate("Boat")
	// if err != nil {
	// 	log.Println("error populate trip store")
	// }
	return nil
}
func (ts *mongoDBTripRepository) Delete(ctx context.Context, id string) error {
	// trip, err := ts.GetByID(ctx, id)
	// if err != nil {
	// 	return err
	// }
	// err = trip.Delete()
	// if err != nil {
	// 	log.Fatal("error in delete transaction ")
	// 	return err
	// }
	return nil
}

func (ts *mongoDBTripRepository) AddPassengers(ctx context.Context, trip *models.Trip, passengers []*models.User) (*models.Trip, error) {
	// Trip := ts.Conn.Model("Trip")
	// Trip.New(trip)
	// trip.Passengers = passengers
	// log.Println("add passenger")
	return nil, nil
}
