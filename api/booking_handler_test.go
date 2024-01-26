package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/DenisBytes/GoHotel/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	db := setUp(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.Store, "Bar Hotel", "italy", 4, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 2, hotel.ID)
	booking := fixtures.AddBooking(db.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0,0,2))
	fmt.Println(booking)
}