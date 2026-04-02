package repositories

import (
	"context"
	"fmt"

	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoCommentRepository implements CommentRepository using MongoDB.
type MongoCommentRepository struct {
	collection *mongo.Collection
}

// NewMongoCommentRepository creates a new MongoCommentRepository.
func NewMongoCommentRepository(collection *mongo.Collection) *MongoCommentRepository {
	return &MongoCommentRepository{collection: collection}
}

// GetAll retrieves all comments from the database.
func (r *MongoCommentRepository) GetAll(ctx context.Context) ([]models.Comment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("failed to decode comments: %w", err)
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	return comments, nil
}

// GetByID retrieves a comment by its ID.
func (r *MongoCommentRepository) GetByID(ctx context.Context, id bson.ObjectID) (*models.Comment, error) {
	var comment models.Comment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("failed to query comment: %w", err)
	}
	return &comment, nil
}

// GetByArticleID retrieves all comments for a specific article.
func (r *MongoCommentRepository) GetByArticleID(ctx context.Context, articleID bson.ObjectID) ([]models.Comment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"article_id": articleID})
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("failed to decode comments: %w", err)
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	return comments, nil
}
