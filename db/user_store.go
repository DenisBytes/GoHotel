package db

import (
	"context"

	"github.com/DenisBytes/GoHotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
}

type MongouserStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongouserStore{

	return &MongouserStore{
		client: client,
		coll : client.Database(DBNAME).Collection(userCollection),

	}
}

func (s *MongouserStore) GetUserByID(ctx context.Context, id string) (*types.User, error){

	//retrieve/convert ObjectID type (for mongodb _id object) from the id string
	oid, err := primitive.ObjectIDFromHex(id)
	if err !=nil {
		return nil,err
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err!=nil{
		return nil, err
	}
	return &user, nil
}