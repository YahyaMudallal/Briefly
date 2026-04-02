package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// User represents a user account in the database.
type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string        `bson:"email" json:"email"`
	FirstName string        `bson:"first_name" json:"firstName"`
	LastName  string        `bson:"last_name" json:"lastName"`
	Password  string        `bson:"password" json:"-"`
	CreatedAt time.Time     `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updatedAt"`
}