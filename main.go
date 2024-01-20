package main

import (
	"context"
	"flag"
	"log"

	"github.com/DenisBytes/GoHotel/api"
	"github.com/DenisBytes/GoHotel/api/middleware"
	"github.com/DenisBytes/GoHotel/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// when something errs, it going to return a json with key "err:" and the value the actual error
var config = fiber.Config{
	ErrorHandler: func (c *fiber.Ctx, err error) error{
		return c.JSON(map[string]string{"error": err.Error()})
	},
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
		store = &db.Store{
			Hotel: hotelStore,
			Room: roomStore,
			User: userStore,
		}
		userHandler = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		authHandler = api.NewAuthHandler(userStore)
		app = fiber.New(config)
		// this is like requestmapping in spring above the controller class. to create a prefixed path.
		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication)
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
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms	)

	//this needs to be at the end
	app.Listen(*listenAddr)

}