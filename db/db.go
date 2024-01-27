package db

import (
	"os"
)

var (
	DBNAME string
	DBURL string
)

func init() {
	// use  ("../.env") for tests. in other location files it can't read the .env file
	// if err := godotenv.Load("../.env"); err != nil {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal(err)
	// }
	DBNAME = os.Getenv("MONGO_DB_NAME")
	DBURL = os.Getenv("MONGO_DB_URL")
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

// this is to make the interface implementable for both mongo and sql.
// BSON == map[string] any      under the hood
type Map map[string]any

type Pagination struct {
	Page  int64
	Limit int64
}


