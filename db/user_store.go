package db

import (
	"context"
	"fmt"

	"github.com/DenisBytes/GoHotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//repo-service

const userCollection = "users"

type Dropper interface{
	Drop(context.Context) error
}
//Interface for multiple Stores/DBs
type UserStore interface {
	Dropper

	GetUserByEmail(context.Context, string) (*types.User, error) 
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error
}

//struct for mongodb store
type MongoUserStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

// "contructor"
func NewMongoUserStore(client *mongo.Client) *MongoUserStore{
	return &MongoUserStore{
		client: client,
		coll : client.Database(DBNAME).Collection(userCollection),
	}
}

//implementation of interface
func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error){
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

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error){
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err!=nil{
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error){
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

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error){
	res, err := s.coll.InsertOne(ctx, user)
	if err!=nil{
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error{
	oid, err := primitive.ObjectIDFromHex(id)
	if err !=nil {
		return err 
	}
	// TODO: it's good to handle if delete didn't delete any user
	// maybe log it or something?

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err!=nil{
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error{
	//values := bson.M{} //bson.M is a map
	values := params.ToBSON()
	update := bson.D{
		{
			Key: "$set", Value: values,
		},
	}
	
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err !=nil {
		return err 
	}
	return nil
}

	func (s *MongoUserStore) Drop(ctx context.Context) error{
		fmt.Println("--- dropping user collection")
		return s.coll.Drop(ctx)
	}