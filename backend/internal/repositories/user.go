package repositories

import (
	"context"
	"fmt"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongoUserRepository implements UserRepository using MongoDB.
type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository creates a new MongoUserRepository.
func NewMongoUserRepository(collection *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{collection: collection}
}

// GetAll retrieves all users from the database.
func (r *MongoUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("%w : failed to query users: %w", apperrors.ErrInternal, err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("%w : failed to decode users: %w", apperrors.ErrInternal, err)
	}

	if users == nil {
		users = []models.User{}
	}

	return users, nil
}

// GetByID retrieves a user by their ID.
func (r *MongoUserRepository) GetByID(ctx context.Context, id bson.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w : user not found", apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w : failed to query user: %w", apperrors.ErrInternal, err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email address.
func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%w : user not found", apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("%w : failed to query user: %w", apperrors.ErrInternal, err)
	}
	return &user, nil
}

// Create creates a new user.
// Returns the model of the created user with the ID field added, or an error
func (r *MongoUserRepository) Create(ctx context.Context, newUser *models.User) (*models.User, error) {
	
	// insert the new user into the database
	result, err := r.collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("%w : failed to create user: %w", apperrors.ErrInternal, err)
	}

	// set the ID of the new user to the inserted ID returned by MongoDB
	userID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, fmt.Errorf("%w : failed to convert inserted ID to ObjectID", apperrors.ErrInternal)
	}
	newUser.ID = userID

	return newUser, nil
}

// Delete remove a user from its ID.
func (r *MongoUserRepository) Delete(ctx context.Context, id bson.ObjectID ) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("%w: failed to delete user: %w", apperrors.ErrInternal, err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("%w: user not found", apperrors.ErrNotFound)
	}
	return nil
}
