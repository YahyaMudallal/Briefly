package repositories

import (
	"context"
	"fmt"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
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
		return nil, fmt.Errorf("%w: failed to query comments: %w", apperrors.ErrInternal, err)
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("%w: failed to decode comments: %w", apperrors.ErrInternal, err)
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
			return nil, fmt.Errorf("%w: comment not found", apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w : failed to query comment: %w", apperrors.ErrInternal, err)
	}

	// comment not found if content is empty 
	if comment.Content == "" {
		return nil, fmt.Errorf("%w: comment not found", apperrors.ErrNotFound)
	}

	return &comment, nil
}

// GetByArticleID retrieves all comments for a specific article.
func (r *MongoCommentRepository) GetByArticleID(ctx context.Context, articleID bson.ObjectID) ([]models.Comment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"article_id": articleID})
	if err != nil {
		return nil, fmt.Errorf("%w: failed to query comments: %w", apperrors.ErrNotFound, err)
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("%w: failed to decode comments: %w", apperrors.ErrInternal, err)
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	return comments, nil
}

// Create inserts a new comment into the database.
func (r *MongoCommentRepository) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	result, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to insert comment: %w", apperrors.ErrInternal, err)
	}

	comment.ID = result.InsertedID.(bson.ObjectID)
	return comment, nil
}

// Delete removes a comment from the database by its ID.
func (r *MongoCommentRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("%w: failed to delete comment: %w", apperrors.ErrInternal, err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("%w: comment not found", apperrors.ErrNotFound)
	}
	return nil
}