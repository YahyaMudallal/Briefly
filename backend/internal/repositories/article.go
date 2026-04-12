package repositories

import (
	"context"
	"fmt"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoArticleRepository implements ArticleRepository using MongoDB.
type MongoArticleRepository struct {
	collection *mongo.Collection
}

// NewMongoArticleRepository creates a new MongoArticleRepository.
func NewMongoArticleRepository(collection *mongo.Collection) *MongoArticleRepository {
	return &MongoArticleRepository{collection: collection}
}

// GetAll retrieves all articles from the database.
func (r *MongoArticleRepository) GetAll(ctx context.Context) ([]models.Article, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("%w : failed to query articles: %w", apperrors.ErrInternal, err)
	}
	defer cursor.Close(ctx)

	var articles []models.Article
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, fmt.Errorf("%w : failed to decode articles: %w", apperrors.ErrInternal, err)
	}

	if articles == nil {
		articles = []models.Article{}
	}

	return articles, nil
}

// GetByID retrieves an article by its ID.
func (r *MongoArticleRepository) GetByID(ctx context.Context, id bson.ObjectID) (*models.Article, error) {
	var article models.Article
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w : article not found", apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w : failed to query article: %w", apperrors.ErrInternal, err)
	}
	return &article, nil
}

// Create inserts a new article into the database.
func (r * MongoArticleRepository) Create(ctx context.Context, article *models.Article) (*models.Article, error) {
	result, err := r.collection.InsertOne(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("%w : failed to insert article: %w", apperrors.ErrInternal, err)
	}
	article.ID = result.InsertedID.(bson.ObjectID)
	return article, nil
}

// Delete removes an article by its ID.
func (r *MongoArticleRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("%w : failed to delete article: %w", apperrors.ErrInternal, err)
	}

	// check if an article was actually deleted
	if result.DeletedCount == 0 {
		return fmt.Errorf("%w : article not found", apperrors.ErrNotFound)
	}

	return nil
}