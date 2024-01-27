package db

const (
	DBNAME = "go-hotel"
	TestDBNAME = "go-hotel-test"
	DBURI = "mongodb://localhost:27017"
)

type Store struct{
	User UserStore
	Hotel HotelStore
	Room RoomStore
	Booking BookingStore
}

//this is to make the interface implementable for both mongo and sql.
//BSON == map[string] any      under the hood
type Map map[string]any


type Pagination struct {
	Page int64
	Limit int64
}