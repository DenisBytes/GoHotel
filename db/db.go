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
}