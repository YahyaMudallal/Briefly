package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Vote represents a user's vote on an article in the database.
type Vote struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ArticleID bson.ObjectID `bson:"article_id" json:"articleId"`
	UserID    bson.ObjectID `bson:"user_id" json:"userId"`
	Type      int           `bson:"type" json:"type"` // upvote=1 or downvote=-1s
}
