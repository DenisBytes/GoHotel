package api

import (
	"github.com/DenisBytes/GoHotel/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler{
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error{
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error{
	return nil
}