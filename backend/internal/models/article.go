package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Article represents a news article in the database.
type Article struct {
	ID            bson.ObjectID 	`bson:"_id,omitempty" json:"id,omitempty"`
	Title         string        	`bson:"title" json:"title"`
	Description   string        	`bson:"description" json:"description"`
	Content       string        	`bson:"content" json:"content"`
	OriginalURL   string        	`bson:"original_url" json:"originalUrl"`
	ImageURL	  string        	`bson:"image_url" json:"imageUrl"`
	Source        string        	`bson:"source" json:"source"`
	UpVotes       int           	`bson:"up_votes" json:"upvotes"`
	DownVotes     int           	`bson:"down_votes" json:"downvotes"`
	Summary       string        	`bson:"summary" json:"summary"`
	CreatedAt     time.Time     	`bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time     	`bson:"updated_at" json:"updatedAt"`
}