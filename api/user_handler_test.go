package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/DenisBytes/GoHotel/db"
	"github.com/DenisBytes/GoHotel/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "go-hotel"

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T){
	if err := tdb.UserStore.Drop(context.TODO()); err!=nil{
		t.Fatal(err)
	}
}

func setUp(t *testing.T) *testdb{
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err!=nil{
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}
func TestPostUser(t *testing.T){
	tdb := setUp(t)
	defer tdb.teardown(t)

	app:=fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email: "some@foo.com",
		FirstName: "James",
		LastName: "Foo",
		Password: "akbfahjbf",
	}
	//Marshal the params to bytes
	b,_ := json.Marshal(params)

	// 3rd param needs an io.Reader (which is an interface). NewReader returns a NewReader form the bytes
	req:= httptest.NewRequest("POST", "/", bytes.NewReader(b))
	//without this it won't know the format/type that we want
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err!=nil {
		t.Log(err)
	}
	// bb, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(bb))
	var user types.User
	//decode the body into the user
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0{
		t.Logf("Expecting a user id")
	}
	if len(user.EncryptedPassword) > 0{
		t.Logf("Expexted the EncryptedPassword not to be included in the json response")
	}
	if user.FirstName != params.FirstName{
		t.Logf("Expected firstName %s but got %s\n", user.FirstName, params.FirstName)
	}
	if user.LastName != params.LastName{
		t.Logf("Expected lastName %s but got %s\n", user.LastName, params.LastName)
	}
	if user.Email != params.Email{
		t.Logf("Expected Email %s but got %s\n", user.Email, params.Email)
	}
}