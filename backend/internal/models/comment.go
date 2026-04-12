package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Comment represents a comment on an article in the database.
type Comment struct {
	ID        bson.ObjectID 	`bson:"_id,omitempty" json:"id,omitempty"`
	AuthorID  bson.ObjectID 	`bson:"author_id" json:"authorId"`
	ArticleID bson.ObjectID 	`bson:"article_id" json:"articleId"`
	Content   string        	`bson:"content" json:"content"`
	ParentID  *bson.ObjectID 	`bson:"parent_id,omitempty" json:"parentId,omitempty"`
	CreatedAt time.Time     	`bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time     	`bson:"updated_at" json:"updatedAt"`
}