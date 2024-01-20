package db

import (
	"context"

	"github.com/DenisBytes/GoHotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// INTERFACE

type RoomStore interface {
	CreateRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

// CLASS AND CONSTRUCTOR
type MongoRoomStore struct {
	client *mongo.Client
	coll *mongo.Collection

	HotelStore
}

func NewMongoRoomStore (client *mongo.Client, hotelStore HotelStore) *MongoRoomStore{
	return &MongoRoomStore{
		client: client,
		coll: client.Database(DBNAME).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

// METHODS

func (s *MongoRoomStore) CreateRoom(ctx context.Context, room *types.Room) (*types.Room, error){
	resp, err := s.coll.InsertOne(ctx, room)
	if err!=nil{
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)

	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := s.HotelStore.UpdateHotel(ctx, filter, update); err!=nil{
		return nil, err
	}
	
	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error){
	res, err := s.coll.Find(ctx, filter)
	if err!=nil{
		return nil , err
	}
	var rooms []*types.Room
	if err:= res.All(ctx,&rooms); err!=nil{
		return nil,err
	}
	return rooms, nil
}