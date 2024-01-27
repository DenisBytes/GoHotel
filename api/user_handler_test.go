package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/DenisBytes/GoHotel/types"
	"github.com/gofiber/fiber/v2"
)

func TestPostUser(t *testing.T){
	tdb := setUp(t)
	defer tdb.teardown(t)

	app:=fiber.New()
	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)

	//Testing post method
	params := types.CreateUserParams{
		Email: "james@foo.com",
		FirstName: "james",
		LastName: "foo",
		Password: "james_foo",
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