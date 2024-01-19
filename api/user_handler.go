package api

import (
	"context"

	"github.com/DenisBytes/GoHotel/db"
	"github.com/DenisBytes/GoHotel/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{
	userStore db.UserStore //not mongoStroe because you may have to access multiple dbs
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "At the watercooler",
	}
	//c.JSON is going to automatically marshal/convert your user into JSON
	return c.JSON(u)
}

func NewUserHandler(userStore db.UserStore) *UserHandler{
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()
	
	user, err := h.userStore.GetUserByID(ctx, id)
	if err!= nil{
		return err
	}
	return c.JSON(user)
}