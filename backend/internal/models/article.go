package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Article represents a news article in the database.
type Article struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Summary     string        `bson:"summary" json:"summary"`
	Content     string        `bson:"content" json:"content"`
	UpVotes     int           `bson:"up_votes" json:"upvotes"`
	DownVotes   int           `bson:"down_votes" json:"downvotes"`
	NbComments  int           `bson:"nb_comments" json:"nbComments"`
	Source      string        `bson:"source" json:"source"`
	CreatedAt   time.Time     `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updatedAt"`
}
