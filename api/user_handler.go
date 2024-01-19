package api

import (
	"errors"

	"github.com/DenisBytes/GoHotel/db"
	"github.com/DenisBytes/GoHotel/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// handlers = return functions of paths

type UserHandler struct{
	userStore db.UserStore //not mongoStore because you may have to access multiple dbs
}

//Constructor for handler
func NewUserHandler(userStore db.UserStore) *UserHandler{
	return &UserHandler{
		userStore: userStore,
	}
}

// Methods. various handlers functions
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err!=nil{
		return err
	}
	//c.JSON is going to automatically marshal/convert your user into JSON
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")	
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err!= nil{
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}
	
	return c.JSON(user)
}


func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error{
	var params types.CreateUserParams
	if err:= c.BodyParser(&params); err!=nil{
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user,err := types.NewUserFromParams(params)
	if err!=nil{
		return err
	}
	insertedUser, err:= h.userStore.CreateUser(c.Context(), user)
	if err!=nil{
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error{
	var (
		//values bson.M
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	if err := c.BodyParser(&params); err!=nil{
		return err
	}

	oid, err := primitive.ObjectIDFromHex(userID)
	if err !=nil {
		return err 
	}
	filter:= bson.M{"_id": oid}
	
	if err:= h.userStore.UpdateUser(c.Context(), filter, params); err!=nil{
		return err
	}
	return c.JSON(map[string]string{"updated": userID})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error{
	userID := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userID); err!=nil{
		return err
	}
	return c.JSON(map[string]string{"deleted": userID})
}