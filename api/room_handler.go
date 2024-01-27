package api

import (
	"context"
	"fmt"
	"time"

	"github.com/DenisBytes/GoHotel/db"
	"github.com/DenisBytes/GoHotel/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CLASS AND CONSTRUCTOR
type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler{
	return &RoomHandler{
		store: store,
	}
}

type BookRoomParams struct {
	NumPersons int `json:"numPersons"`
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
}

func (p *BookRoomParams) validate() error{

	now:= time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate){
		return fmt.Errorf("Cannot book a room in the past")
	}

	return nil
}

func (h *RoomHandler) HandleBookRoom (c *fiber.Ctx) error{

	var params BookRoomParams
	if err := c.BodyParser(&params); err !=nil {
		return err
	}

	//validate the body params
	if err := params.validate(); err!=nil{
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg: "internal server error",
		})
	}

	ok, err = h.isRoomAvailableForBooking(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg: fmt.Sprintf("room %s already booked", c.Params("id")),
		})
	}

	booking := types.Booking{
		UserID: user.ID,
		RoomID: roomID,
		FromDate: params.FromDate,
		TillDate: params.TillDate,
		NumPersons: params.NumPersons,
	}

	inserted, err := h.store.Booking.CreateBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context ,roomID primitive.ObjectID, params BookRoomParams ) (bool, error){

	where := db.Map{
        "roomID": roomID,
        "$or": []bson.M{
            // Check if the existing booking starts during the new booking period
            {
                "fromDate": bson.M{"$gte": params.FromDate, "$lte": params.TillDate},
            },
            // Check if the existing booking ends during the new booking period
            {
                "tillDate": bson.M{"$gte": params.FromDate, "$lte": params.TillDate},
            },
            // Check if the existing booking completely contains the new booking period
            {
                "fromDate": bson.M{"$lte": params.FromDate},
                "tillDate": bson.M{"$gte": params.TillDate},
            },
        },
    }
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0
	return ok, nil
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), db.Map{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}