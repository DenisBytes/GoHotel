package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
	minFirstNameLen = 2
	minLastNameLen = 2
	minPasswordLen = 8
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // to omit value on respond/render. if omit only the empty = omitempty. if always omit ? json:"-"
	FirstName string `bson:"firstName" json:"firstName"`
	LastName string `bson:"lastName" json:"lastName"`
	Email string `bson:"email" json:"email"`
	EncryptedPassword string `bson:"encryptedPassword" json:"-"`
	IsAdmin bool `bson:"isAdmin" json:"isAdmin"`
}

//DTO
type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

//DTO
type UpdateUserParams struct{
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}


//DTO Constructor
func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err!= nil{
		return nil, err
	}
	return &User{
		FirstName: params.FirstName,
		LastName: params.LastName,
		Email: params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}



func (params CreateUserParams) Validate() map[string]string{
	errors := map[string]string{} 

	if len(params.FirstName) < minFirstNameLen{
		errors["firstName"] = fmt.Sprintf("First name should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen{
		errors["lastName"] = fmt.Sprintf("Last name should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen{
		errors["password"] = fmt.Sprintf("Password should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email){
		errors["email"] = fmt.Sprintf("Email %s is invalid", params.Email)
	}
	return errors
}

func isEmailValid(e string) bool {
	// not very good regexo. is by-passable
    emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
    return emailRegex.MatchString(e)
}

func IsValidPassword(encpw, pw string) bool{
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

func (p UpdateUserParams) ToBSON() bson.M{
	m := bson.M{}
	if len(p.FirstName) > 0{
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0{
		m["lastName"] = p.LastName
	}
	return m
}