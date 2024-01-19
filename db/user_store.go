package db

import (
	"context"

	"github.com/DenisBytes/GoHotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
//repo-service

const userCollection = "users"

//Interface for multiple Stores/DBs
type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
}

//struct for mongodb store
type MongouserStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

// "contructor"
func NewMongoUserStore(client *mongo.Client) *MongouserStore{
	return &MongouserStore{
		client: client,
		coll : client.Database(DBNAME).Collection(userCollection),
	}
}

//implementation of interface
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

func (s *MongouserStore) GetUsers(ctx context.Context) ([]*types.User, error){
	cur, err := s.coll.Find(ctx, bson.M{})
	if err!=nil{
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err!=nil{
		return nil,err
	}

	return users, nil
}

func (s *MongouserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error){
	res, err := s.coll.InsertOne(ctx, user)
	if err!=nil{
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}