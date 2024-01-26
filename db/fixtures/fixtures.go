package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DenisBytes/GoHotel/db"
	"github.com/DenisBytes/GoHotel/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fname, lname string, admin bool) *types.User{

	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: fmt.Sprintf("%s@%s.com", fname, lname),
		FirstName: fname,
		LastName: lname,
		Password: fmt.Sprintf("%s_%s", fname, lname),
	})
	if err!=nil{
		log.Fatal(err)
	}
	user.IsAdmin = admin

	insertedUser, err := store.User.CreateUser(context.TODO(), user)
	if err!=nil{
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store *db.Store, name string, location string, rating int, rooms []primitive.ObjectID) *types.Hotel{
	var roomsIDS =  rooms
	if rooms == nil {
		roomsIDS = []primitive.ObjectID{}
	}

	hotel := types.Hotel{
		Name: name,
		Location: location,
		Rooms: roomsIDS,
		Rating: rating,
	}
	insertedHotel, err := store.Hotel.CreateHotel(context.TODO(), &hotel)
	if err!=nil{
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room{
	room := &types.Room{
		Size: size,
		SeaSide: seaside,
		Price: price,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Room.CreateRoom(context.Background(), room)
	if err!=nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID: userID,
		RoomID: roomID,
		FromDate: from,
		TillDate: till,
	}

	insertedBooking ,err:= store.Booking.CreateBooking(context.Background(), booking);
	if  err!=nil{
		log.Fatal(err)
	}

	return insertedBooking
}