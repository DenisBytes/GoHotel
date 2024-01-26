package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DenisBytes/GoHotel/api"
	"github.com/DenisBytes/GoHotel/db"
	"github.com/DenisBytes/GoHotel/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err!=nil{
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err!=nil{
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	store := &db.Store{
		User: db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room: db.NewMongoRoomStore(client, hotelStore),
		Hotel: hotelStore,
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Println("Not-Admin: ", api.CreateTokenFromUser(user))

	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("Admin", api.CreateTokenFromUser(admin))

	hotel := fixtures.AddHotel(store, "Bellucia", "France", 4, nil)
	fmt.Println(hotel)

	room := fixtures.AddRoom(store, "large", true, 88.44, hotel.ID)
	fmt.Println(room)

	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0,0,2))
	fmt.Println(booking) 
}