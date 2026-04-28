package repositories

import (
	"context"

	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ArticleRepository defines methods for article data access.
type ArticleRepository interface {
	GetAll(ctx context.Context) ([]models.Article, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*models.Article, error)
	Create(ctx context.Context, article *models.Article) (*models.Article, error)
	Delete(ctx context.Context, id bson.ObjectID) error
	IncrementCommentCount(ctx context.Context, articleID bson.ObjectID, amount int) error
	IncrementUpVotes(ctx context.Context, articleID bson.ObjectID, amount int) error
	IncrementDownVotes(ctx context.Context, articleID bson.ObjectID, amount int) error
}

// UserRepository defines methods for user data access.
type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id bson.ObjectID) error
}

// CommentRepository defines methods for comment data access.
type CommentRepository interface {
	GetAll(ctx context.Context) ([]models.Comment, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*models.Comment, error)
	GetByArticleID(ctx context.Context, articleID bson.ObjectID) ([]models.Comment, error)
	Create(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	Delete(ctx context.Context, id bson.ObjectID) error
	Update(ctx context.Context, comment *models.Comment) error
	DeleteByArticleID(ctx context.Context, articleID bson.ObjectID) error
}

// VoteRepository defines methods for vote data access.
type VoteRepository interface {
	GetByArticleIDAndUserID(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) (*models.Vote, bool, error)
	GetAllByArticleID(ctx context.Context, articleID bson.ObjectID) ([]models.Vote, []models.Vote, error)
	GetAllByUserID(ctx context.Context, userID bson.ObjectID) ([]models.Vote, error)
	ToggleUpvote(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) (int, error)
	ToggleDownvote(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) (int, error)
}
