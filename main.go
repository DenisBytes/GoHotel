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

const dburi = "mongodb://localhost:27017"
const dbname = "go-hotel"
const userColl = "users"

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

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err!=nil{
		log.Fatal(err)
	}
	// handlers initialiazation
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	// this is like requestmapping in spring above the controller class. to create a prefixed path.
	apiv1 := app.Group("/api/v1")
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	app.Listen(*listenAddr)

}