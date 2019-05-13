package repository

import (
	"context"
	"log"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/user"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type mongoDBUserRepository struct {
	Conn *mongodm.Connection
}

func NewMongoDBUserRepository(Conn *mongodm.Connection) user.Repository {
	return &mongoDBUserRepository{Conn}
}

func (m *mongoDBUserRepository) fetch(ctx context.Context, query interface{}, limit int, skip int, sort string) (res []*models.User, nextSkip int, err error) {
	User := m.Conn.Model("User")
	user := []*models.User{}
	err = User.Find(query).Sort(sort).Skip(skip).Limit(limit).Populate("TripHistory").Exec(&user)
	if err != nil {
		log.Fatal("error in find repo")
		return nil, skip + limit, err
	}
	return user, skip + limit, nil
}

func (m *mongoDBUserRepository) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.User, nextSkip int, err error) {
	query := bson.M{"deleted": false}
	if sort == "" {
		sort = "_modifiedAt"
	}
	result, nextSkip, err := m.fetch(ctx, query, limit, skip, sort)
	if err != nil {
		log.Fatal("error in fetch user ")
		return nil, nextSkip, err
	}
	return result, nextSkip, nil
}

func (m *mongoDBUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	User := m.Conn.Model("User")
	user := &models.User{}
	err := User.FindId(bson.ObjectIdHex(id)).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		// log.Fatal("user not found")
		return nil, err
	} else if err != nil {
		// log.Fatal("database error")
		return nil, err
	}
	return user, nil
}
func (m *mongoDBUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	User := m.Conn.Model("User")
	user := &models.User{}
	err := User.FindOne(bson.M{"name": username, "deleted": false}).Exec(user)
	if err != nil {
		log.Fatal("error find username")
		return nil, err
	}
	return user, nil
}
func (m *mongoDBUserRepository) Update(ctx context.Context, selector interface{}, update interface{}) error {
	// User := m.Conn.Model("User")
	// err := User.Update(selector, update)
	// if err != nil {
	// 	log.Fatal("error " + err.Error())
	// 	return err
	// }
	return nil
}
func (m *mongoDBUserRepository) Store(ctx context.Context, user *models.User) error {
	// User := m.Conn.Model("User")
	// User.New(user)
	// err := user.Save()
	// if err != nil {
	// 	fmt.Println("error " + err.Error())
	// 	return err
	// }

	return nil
}
func (m *mongoDBUserRepository) Delete(ctx context.Context, id string) error {
	// user, err := m.GetByID(ctx, id)
	// if err != nil {
	// 	log.Fatal("error when get user by id")
	// 	return err
	// }
	// err = user.Delete()
	// if err != nil {
	// 	log.Fatal("error when delete post")
	// 	return err
	// }
	return nil

}
