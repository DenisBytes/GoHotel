package db

import (
	"context"

	"github.com/DenisBytes/GoHotel/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// INTERFACE FOR DB
type HotelStore interface {
	CreateHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, Map, Map) error
	GetHotels(context.Context, Map, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

// CLASS AND CONSTRUCTOR
type MongoHotelStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoHotelStore (client *mongo.Client) *MongoHotelStore{
	return &MongoHotelStore{
		client: client,
		coll: client.Database(DBNAME).Collection("hotels"),
	}
}

// METHODS

func (s *MongoHotelStore) CreateHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error){
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err!=nil{
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context,filter Map, pagFilter *Pagination ) ([]*types.Hotel, error){
	opts := options.FindOptions{}
	opts.SetSkip((pagFilter.Page -1)  * pagFilter.Limit)
	opts.SetLimit(pagFilter.Limit)
	resp, err := s.coll.Find(ctx, filter, &opts)
	if err!=nil{
		return nil, err
	}
	var hotels []*types.Hotel
	if err:= resp.All(ctx, &hotels); err!=nil{
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter Map, update Map) error{
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error){
	oid, err := primitive.ObjectIDFromHex(id)
	if err !=nil {
		return nil, err
	}
	var hotel *types.Hotel
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err!=nil{
		return nil, err
	}

	return hotel, nil
}