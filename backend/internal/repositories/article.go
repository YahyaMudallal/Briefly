package repositories

import (
	"context"
	"fmt"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

// GetPaginated retrieves a paginated list of articles from the database.
func (r *MongoArticleRepository) GetPaginated(ctx context.Context, page int, limit int, sortBy string, order string) ([]models.Article, error) {
	skip := int64((page - 1) * limit)
	limit64 := int64(limit)

	// determine sort direction
	sortDir := -1
    if order == "asc" {
        sortDir = 1
    }

	pipeline := mongo.Pipeline{}

	if sortBy == "hotness" {
		// HOTNESS ALGORITHM
		// Formula : (Upvotes - Downvotes) / (Age in hours + 1)

		// numerator : up_votes - down_votes
		numerator := bson.M{"$subtract": bson.A{"$up_votes", "$down_votes"}}

		// compute the age : $$NOW (current date) - published_at = result in milliseconds
		ageInMillis := bson.M{"$subtract": bson.A{"$$NOW", "$published_at"}}

		// convert milliseconds to hours (1000 * 60 * 60 = 3600000)
		ageInHours := bson.M{"$divide": bson.A{ageInMillis, 3600000}}

		// denominator : Age in hours + 1 (to avoid division by zero)
		denominator := bson.M{"$add": bson.A{ageInHours, 1}}

		// final score
		hotnessScore := bson.M{"$divide": bson.A{numerator, denominator}}

		// inject the hotness score into the pipeline
		pipeline = append(pipeline, bson.D{{Key: "$addFields", Value: bson.M{"hotness_score": hotnessScore}}})

		// sort by hotness score
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "hotness_score", Value: sortDir}, {Key: "_id", Value: 1}}}})

	} else {
		// sort by date
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "published_at", Value: sortDir}, {Key: "_id", Value: 1}}}})
	}

	// pagination
	pipeline = append(pipeline, bson.D{{Key: "$skip", Value: skip}})
	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: limit64}})
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("%w : failed to aggregate articles: %w", apperrors.ErrInternal, err)
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
func (r *MongoArticleRepository) Create(ctx context.Context, article *models.Article) (*models.Article, error) {
	result, err := r.collection.InsertOne(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("%w : failed to insert article: %w", apperrors.ErrInternal, err)
	}
	article.ID = result.InsertedID.(bson.ObjectID)
	return article, nil
}

// CreateMany inserts new articles into the database.
func (r *MongoArticleRepository) CreateMany(ctx context.Context, articles []models.Article) ([]models.Article, error) {

	// convert the typed array in an interface array
	var docs []interface{}
	for _, article := range articles {
		docs = append(docs, article)
	} 

	opts := options.InsertMany().SetOrdered(false)

	// inset the interface array
	result, err := r.collection.InsertMany(ctx, docs, opts)
	if err != nil {
		return nil, fmt.Errorf("%w : failed to insert articles: %w", apperrors.ErrInternal, err )
	}

	// assign new IDs
	for i, id := range result.InsertedIDs {
		articles[i].ID = id.(bson.ObjectID)
	}

	return articles, nil
}

// Update updates an existing article in the database.
func (r *MongoArticleRepository) Update(ctx context.Context, article *models.Article) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": article.ID}, bson.M{"$set": article})
	if err != nil {
		return fmt.Errorf("%w : failed to update article: %w", apperrors.ErrInternal, err)
	}
	return nil
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

// IncrementCommentCount securely updates the comment count using atomic $inc
func (r *MongoArticleRepository) IncrementCommentCount(ctx context.Context, articleID bson.ObjectID, amount int) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": articleID},
		bson.M{"$inc": bson.M{"nb_comments": amount}},
	)
	if err != nil {
		return fmt.Errorf("%w: failed to increment article comment count: %w", apperrors.ErrInternal, err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%w: article not found", apperrors.ErrNotFound)
	}
	return nil
}

// IncrementUpVotes securely updates the upvote count using atomic $inc
func (r *MongoArticleRepository) IncrementUpVotes(ctx context.Context, articleID bson.ObjectID, amount int) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": articleID},
		bson.M{"$inc": bson.M{"up_votes": amount}},
	)
	if err != nil {
		return fmt.Errorf("%w: failed to increment article upvote count: %w", apperrors.ErrInternal, err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%w: article not found", apperrors.ErrNotFound)
	}
	return nil
}

// IncrementDownVotes securely updates the downvote count using atomic $inc
func (r *MongoArticleRepository) IncrementDownVotes(ctx context.Context, articleID bson.ObjectID, amount int) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": articleID},
		bson.M{"$inc": bson.M{"down_votes": amount}},
	)
	if err != nil {
		return fmt.Errorf("%w: failed to increment article downvote count: %w", apperrors.ErrInternal, err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%w: article not found", apperrors.ErrNotFound)
	}
	return nil
}
