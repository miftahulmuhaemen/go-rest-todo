package user

import (
	"context"
	"go-rest-todo/business"
	core "go-rest-todo/core/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepository struct {
	col *mongo.Collection
}

type Collection struct {
	ID        primitive.ObjectID `bson:"_id"`
	RoleID    primitive.ObjectID `bson:"role_id"`
	Name      string             `bson:"name"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func NewCollection(user core.User) (*Collection, error) {
	var objectID primitive.ObjectID
	if user.ID == "" {
		objectID = primitive.NewObjectID()
	} else {
		userID, err := primitive.ObjectIDFromHex(user.ID)
		if err != nil {
			return nil, business.ErrInvalidID
		}
		objectID = userID
	}

	roleObjectID, err := primitive.ObjectIDFromHex(user.RoleID)
	if err != nil {
		return nil, business.ErrInvalidID
	}

	return &Collection{
		ID:        objectID,
		RoleID:    roleObjectID,
		Name:      user.Name,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.Updatet,
	}, nil
}

func (col *Collection) ToCore() core.User {
	return core.User{
		ID:        col.ID.Hex(),
		RoleID:    col.RoleID.Hex(),
		Name:      col.Name,
		Username:  col.Username,
		Password:  col.Password,
		CreatedAt: col.CreatedAt,
		Updatet:   col.UpdatedAt,
	}
}

func NewMongoDBRepository(db *mongo.Database) *MongoDBRepository {
	return &MongoDBRepository{
		db.Collection("users"),
	}
}

func (repo *MongoDBRepository) Create(user core.User) (core.User, error) {
	col, err := NewCollection(user)
	if err != nil {
		return core.User{}, err
	}

	_, err = repo.col.InsertOne(context.TODO(), col)
	if err != nil {
		return core.User{}, err
	}

	return col.ToCore(), nil
}
