package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/DenisBytes/GoHotel/api/middleware"
	"github.com/DenisBytes/GoHotel/db/fixtures"
	"github.com/DenisBytes/GoHotel/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	db := setUp(t)
	defer db.teardown(t)

	var (
		adminUser = fixtures.AddUser(db.Store, "admin", "admin", true)
		user = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel = fixtures.AddHotel(db.Store, "Bar Hotel", "italy", 4, nil)
		room = fixtures.AddRoom(db.Store, "small", true, 2, hotel.ID)
		booking = fixtures.AddBooking(db.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0,0,2))
		app = fiber.New()
		adminRoute = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)
	_ = booking
	adminRoute.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("x-api-token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK{
		t.Fatalf("expected 200 but got : %d\n", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err!= nil {
		t.Fatal(err)
	}
	
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking but got: %+v\n", bookings)
	}

	//the bookings are not equal. beacause of the timestamps of the dates being  different in millliseconds
	bookings[0].FromDate = booking.FromDate
	bookings[0].TillDate = booking.TillDate
	if !reflect.DeepEqual(booking, bookings[0]){
		fmt.Printf("%+v\n",  booking)
		fmt.Printf("%+v\n",  bookings[0])
		t.Fatal("expected bookings to be equal")
	}

	//testing non admin user
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("x-api-token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK{
		t.Fatalf("expected a non 200 but got : %d\n", resp.StatusCode)
	}

}

func TestUserGetBookingByID(t *testing.T){
	db := setUp(t)
	defer db.teardown(t)
	var (
		userNotBooker = fixtures.AddUser(db.Store, "jimmy", "watercooler", false)
		user = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel = fixtures.AddHotel(db.Store, "Bar Hotel", "italy", 4, nil)
		room = fixtures.AddRoom(db.Store, "small", true, 2, hotel.ID)
		booking = fixtures.AddBooking(db.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0,0,2))
		app = fiber.New()
		route = app.Group("/", middleware.JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("x-api-token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK{
		t.Fatalf("expected 200 but got: %d\n", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	//the bookings are not equal. beacause of the timestamps of the dates being  different in millliseconds
	bookingResp.FromDate = booking.FromDate
	bookingResp.TillDate = booking.TillDate
	if !reflect.DeepEqual(booking, bookingResp){
		fmt.Printf("%+v\n",  booking)
		fmt.Printf("%+v\n",  bookingResp)
		t.Fatal("expected bookings to be equal")
	}


	//testing if another user that didn't book this particular booking can access the route
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("x-api-token", CreateTokenFromUser(userNotBooker))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK{
		t.Fatalf("expected a non 200 but got : %d\n", resp.StatusCode)
	}
}