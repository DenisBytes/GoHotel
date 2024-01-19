package types

type User struct {
	//BSON: It is a binary representation of JSON-like documents, designed to be efficient for storage and data interchange.
	ID        string `bson:"_id,omitempty" json:"id,omitempty"` // to omit value on respond/render. if omit only the empty = omitempty. if always omit ? json:"_"
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
}
