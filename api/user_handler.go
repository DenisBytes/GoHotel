package api

import (
	"github.com/DenisBytes/GoHotel/db"
	"github.com/DenisBytes/GoHotel/types"
	"github.com/gofiber/fiber/v2"
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