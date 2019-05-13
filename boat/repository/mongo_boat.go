package repository

import (
	"context"
	"log"

	"github.com/nolan23/kapaltoba-backend/boat"
	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type mongoDBBoatRepository struct {
	Conn *mongodm.Connection
}

func NewMongoDBBoatRepository(Conn *mongodm.Connection) boat.Repository {
	return &mongoDBBoatRepository{Conn}
}

func (m *mongoDBBoatRepository) fetch(ctx context.Context, query interface{}, limit int, skip int, sort string) (res []*models.Boat, nextSkip int, err error) {
	Boat := m.Conn.Model("Boat")
	boat := []*models.Boat{}
	err = Boat.Find(query).Sort(sort).Skip(skip).Limit(limit).Exec(&boat)
	if err != nil {
		log.Fatal("error in find boat " + err.Error())
		return nil, skip + limit, err
	}
	return boat, skip + limit, nil
}

func (m *mongoDBBoatRepository) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Boat, nextSkip int, err error) {
	query := bson.M{"deleted": false}
	if sort == "" {
		sort = "_modifiedAt"
	}
	result, nextSkip, err := m.fetch(ctx, query, limit, skip, sort)
	if err != nil {
		log.Fatal("error in fetch boat " + err.Error())
		return nil, nextSkip, err
	}
	return result, nextSkip, nil
}

func (m *mongoDBBoatRepository) GetByID(ctx context.Context, id string) (*models.Boat, error) {
	Boat := m.Conn.Model("Boat")
	boat := &models.Boat{}
	err := Boat.FindId(bson.ObjectIdHex(id)).Exec(boat)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return boat, nil
}

func (m *mongoDBBoatRepository) Update(ctx context.Context, selector interface{}, update interface{}) error {
	Boat := m.Conn.Model("Boat")
	err := Boat.Update(selector, update)
	if err != nil {
		log.Fatal("error update boat " + err.Error())
		return err
	}
	return nil
}

func (m *mongoDBBoatRepository) Store(ctx context.Context, boat *models.Boat) error {
	// Boat := m.Conn.Model("Boat")
	// Boat.New(boat)
	// err := boat.Save()
	// if err != nil {
	// 	log.Fatal("error in stre boat " + err.Error())
	// 	return err
	// }
	return nil
}

func (m *mongoDBBoatRepository) Delete(ctx context.Context, id string) error {
	// boat, err := m.GetByID(ctx, id)
	// if err != nil {
	// 	return err
	// }
	// err = boat.Delete()
	// if err != nil {
	// 	log.Fatal("error in delete boat ")
	// 	return err
	// }
	return nil
}
