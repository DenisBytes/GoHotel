package main

import (
	"context"
	"log"

	"github.com/DenisBytes/GoHotel/db"
	"github.com/DenisBytes/GoHotel/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	roomStore db.RoomStore
	hotelStore db.HotelStore
	ctx = context.Background()
)

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Rooms: []primitive.ObjectID{},
		Rating: rating,
	}
	rooms := []types.Room{
		{
			Size: "small",
			Price: 99.9,
		},
		{
			Size: "normal",
			Price: 122.9,
		},
		{
			Size: "kingSize",
			Price: 222.9,
		},
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err!=nil{
		log.Fatal(err)
	}

	for _,room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.CreateRoom(ctx, &room)
		if err!=nil{
			log.Fatal(err)
		}
	}

}

func main() {
	seedHotel("Bellucia", "France", 3)
	seedHotel("The Cozy Hotel", "Netherlands", 4)
	seedHotel("Don't die in your sleep", "London", 1)


}

func init(){
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err!=nil{
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err!=nil{
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}