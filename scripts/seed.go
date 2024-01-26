package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DenisBytes/GoHotel/api"
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
	userStore  db.UserStore
	bookingStore db.BookingStore
	ctx = context.Background()
)

func seedUser(isAdmin bool, email, fname, lname, password string ) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: email,
		FirstName: fname,
		LastName: lname,
		Password: password,
	})
	if err!=nil{
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin

	insertedUser, err := userStore.CreateUser(context.TODO(), user)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))

	return insertedUser
}

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name: name,
		Location: location,
		Rooms: []primitive.ObjectID{},
		Rating: rating,
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err!=nil{
		log.Fatal(err)
	}

	return insertedHotel
}

func seedRoom(size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room{
	room := &types.Room{
		Size: size,
		SeaSide: seaside,
		Price: price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.CreateRoom(context.Background(), room)
	if err!=nil {
		log.Fatal(err)
	}
	return insertedRoom

}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID: userID,
		RoomID: roomID,
		FromDate: from,
		TillDate: till,
	}

	if _,err:= bookingStore.CreateBooking(context.Background(), booking); err!=nil{
		log.Fatal(err)
	}

}


func main() {
	u1 := seedUser(false,"james@foo.com", "james", "foo", "supersecurepassword")
	u2 := seedUser(true,"admin@admin.com", "admin", "admin", "adminpassword123")

	h1 := seedHotel("Bellucia", "France", 3)
	h2 := seedHotel("The Cozy Hotel", "Netherlands", 4)
	h3 := seedHotel("Don't die in your sleep", "London", 1)

	r1 := seedRoom("small", true, 89.99, h1.ID)
	seedRoom("medium", true, 189.99, h2.ID)
	seedRoom("large", true, 289.99, h3.ID)
	r4 := seedRoom("small", true, 89.99, h1.ID)
	seedRoom("medium", true, 189.99, h2.ID)
	seedRoom("large", true, 289.99, h3.ID)
	seedRoom("small", true, 89.99, h1.ID)
	seedRoom("medium", true, 189.99, h2.ID)
	seedRoom("large", true, 289.99, h3.ID)

	seedBooking(u1.ID, r1.ID, time.Now(), time.Now().AddDate(0,0,2))
	seedBooking(u2.ID, r4.ID, time.Now(), time.Now().AddDate(0,0,2))
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
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}