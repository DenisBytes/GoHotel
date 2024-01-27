package main

import (
	"context"
	"flag"
	"log"

	"github.com/DenisBytes/GoHotel/api"
	"github.com/DenisBytes/GoHotel/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// when something errs, it going to return a json with key "err:" and the value the actual error
var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	// this is a flag for command. name of flag - default value - description
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err!=nil{
		log.Fatal(err)
	}

	// handlers initialiazation
	var (
		hotelStore = db.NewMongoHotelStore(client)
		roomStore = db.NewMongoRoomStore(client, hotelStore)
		userStore = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store = &db.Store{
			Hotel: hotelStore,
			Room: roomStore,
			User: userStore,
			Booking: bookingStore,
		}
		userHandler = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		authHandler = api.NewAuthHandler(userStore)
		roomHandler = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app = fiber.New(config)
		apiv1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
		// we use apiv1 and not app becuase jwt is included in apiv1 and in app no
		admin = apiv1.Group("/admin", api.AdminAuth)
		auth = app.Group("/api")
	)

	//Auth handlers
	auth.Post("/auth", authHandler.HandleAuthenticate)

	//user handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	//hotel handlers
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	//rrìoom handlers
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	admin.Get("/bookings", bookingHandler.HandleGetBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	//this needs to be at the end
	app.Listen(*listenAddr)

}