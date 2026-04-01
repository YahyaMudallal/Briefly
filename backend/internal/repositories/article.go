package repositories

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}
	defer cursor.Close(ctx)

	var articles []models.Article
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, fmt.Errorf("failed to decode articles: %w", err)
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
			return nil, fmt.Errorf("article not found")
		}
		return nil, fmt.Errorf("failed to query article: %w", err)
	}
	return &article, nil
}
