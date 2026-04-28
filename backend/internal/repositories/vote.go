package repositories

import (
	"context"
	"fmt"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoVoteRepository implements VoteRepository using MongoDB.

type MongoVoteRepository struct {
	collection *mongo.Collection
}

// NewMongoVoteRepository creates a new MongoVoteRepository.
func NewMongoVoteRepository(collection *mongo.Collection) *MongoVoteRepository {
	return &MongoVoteRepository{collection: collection}
}

// GetByArticleIDAndUserID retrieves a vote by article ID and user ID.
func (r *MongoVoteRepository) GetByArticleIDAndUserID(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) (*models.Vote, bool, error) {
	var vote models.Vote
	err := r.collection.FindOne(ctx, bson.M{"article_id": articleID, "user_id": userID}).Decode(&vote)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, false, nil // No vote found for this article and user; this is not an error, just means the user hasn't voted yet
		}
		// An error occurred while querying the database
		return nil, false, fmt.Errorf("%w : failed to query vote: %w", apperrors.ErrInternal, err)
	}
	return &vote, true, nil
}

// GetAllByArticleID retrieves all votes for a specific article.
func (r *MongoVoteRepository) GetAllByArticleID(ctx context.Context, articleID bson.ObjectID) ([]models.Vote, []models.Vote, error) {
	cursorUp, err := r.collection.Find(ctx, bson.M{"article_id": articleID, "type": 1})
	cursorDown, err := r.collection.Find(ctx, bson.M{"article_id": articleID, "type": -1})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: failed to query votes: %w", apperrors.ErrInternal, err)
	}
	defer cursorUp.Close(ctx)
	defer cursorDown.Close(ctx)

	var upvotes []models.Vote
	var downvotes []models.Vote
	if cursorUp.All(ctx, &upvotes) != nil || cursorDown.All(ctx, &downvotes) != nil {
		return nil, nil, fmt.Errorf("%w: failed to decode votes: %w", apperrors.ErrInternal, err)
	}

	if upvotes == nil {
		upvotes = []models.Vote{}
	}
	if downvotes == nil {
		downvotes = []models.Vote{}
	}

	return upvotes, downvotes, nil
}

// GetAllByUserID retrieves all votes made by a specific user.
func (r *MongoVoteRepository) GetAllByUserID(ctx context.Context, userID bson.ObjectID) ([]models.Vote, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("%w: failed to query votes: %w", apperrors.ErrInternal, err)
	}
	defer cursor.Close(ctx)

	var votes []models.Vote
	if err = cursor.All(ctx, &votes); err != nil {
		return nil, fmt.Errorf("%w: failed to decode votes: %w", apperrors.ErrInternal, err)
	}

	if votes == nil {
		votes = []models.Vote{}
	}

	return votes, nil
}

// ToggleUpvote toggles an upvote for a specific article and user.
func (r *MongoVoteRepository) ToggleUpvote(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) (int, error) {
	//return value: (prev_status, error_status); for ex: (1, nil) if prev status was upvote, (-1, nil) if it was downvote, and (0, nil) if it was novote
	// 0, err if there was an error

	vote, found, err := r.GetByArticleIDAndUserID(ctx, articleID, userID)

	//if error return 0, err
	if err != nil {
		return 0, fmt.Errorf("%w : failed to query vote: %w", apperrors.ErrInternal, err)
	}

	// User has not voted yet, create a new upvote
	// novote ->  upvote
	if !found {
		_, err := r.collection.InsertOne(ctx, models.Vote{
			ArticleID: articleID,
			UserID:    userID,
			Type:      1,
		})
		if err != nil {
			return 0, fmt.Errorf("%w : failed to insert upvote: %w", apperrors.ErrInternal, err)
		}
		return 0, nil
	}

	// User has already downvoted, change to upvote
	// downvote -> upvote
	if vote.Type == -1 {
		_, err = r.collection.UpdateOne(ctx, bson.M{"_id": vote.ID}, bson.M{"$set": bson.M{"type": 1}})
		if err != nil {
			return 0, fmt.Errorf("%w : failed to update vote to upvote: %w", apperrors.ErrInternal, err)
		}

		return -1, nil
	}

	// User has already upvoted, remove the upvote
	// upvote -> novote
	if vote.Type == 1 {
		_, err := r.collection.DeleteOne(ctx, bson.M{"_id": vote.ID})
		if err != nil {
			return 0, fmt.Errorf("%w : failed to delete upvote: %w", apperrors.ErrInternal, err)
		}
		return 1, nil
	}

	// This should never happen, but if it does, return 0 with no error to indicate no change
	return 0, nil
}

// ToggleDownvote toggles a downvote for a specific article and user.
func (r *MongoVoteRepository) ToggleDownvote(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) (int, error) {
	//return value: (prev_status, error_status); for ex: (1, nil) if prev status was upvote, (-1, nil) if it was downvote, and (0, nil) if it was novote
	// 0, err if there was an errors

	vote, found, err := r.GetByArticleIDAndUserID(ctx, articleID, userID)
	if err != nil {
		return 0, fmt.Errorf("%w : failed to query vote: %w", apperrors.ErrInternal, err)
	}

	// User has not voted yet, create a new downvote
	// novote -> downvote
	if !found {
		_, err := r.collection.InsertOne(ctx, models.Vote{
			ArticleID: articleID,
			UserID:    userID,
			Type:      -1,
		})
		if err != nil {
			return 0, fmt.Errorf("%w : failed to insert downvote: %w", apperrors.ErrInternal, err)
		}
		return 0, nil
	}

	// User has already upvoted, change to downvote
	// upvote -> downvote
	if vote.Type == 1 {
		_, err = r.collection.UpdateOne(ctx, bson.M{"_id": vote.ID}, bson.M{"$set": bson.M{"type": -1}})
		if err != nil {
			return 0, fmt.Errorf("%w : failed to update vote to downvote: %w", apperrors.ErrInternal, err)
		}
		return 1, nil
	}

	// User has already downvoted, remove the downvote
	// downvote -> novote
	if vote.Type == -1 {
		_, err := r.collection.DeleteOne(ctx, bson.M{"_id": vote.ID})
		if err != nil {
			return 0, fmt.Errorf("%w : failed to delete downvote: %w", apperrors.ErrInternal, err)
		}
		return -1, nil
	}

	return 0, nil
}
