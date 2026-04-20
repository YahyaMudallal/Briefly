package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
	"github.com/YahyaMudallal/newsWebSite/internal/clients"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"github.com/YahyaMudallal/newsWebSite/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ArticleService provides business logic for articles.
type ArticleService struct {
	articleRepo repositories.ArticleRepository
	userRepo repositories.UserRepository
	commentRepo repositories.CommentRepository
	newsClient clients.NewsClient
}

// NewArticleService creates a new ArticleService.
func NewArticleService(articleRepo repositories.ArticleRepository, userRepo repositories.UserRepository, commentRepo repositories.CommentRepository, newsClient clients.NewsClient) *ArticleService {
	return &ArticleService{articleRepo: articleRepo, userRepo: userRepo, commentRepo: commentRepo, newsClient: newsClient}
}

// GetAllArticles retrieves all articles.
func (s *ArticleService) GetAllArticles(ctx context.Context) ([]models.Article, error) {
	return s.articleRepo.GetAll(ctx)
}

// GetArticleByID retrieves an article by its ID.
func (s *ArticleService) GetArticleByID(ctx context.Context, id bson.ObjectID) (*models.Article, error) {
	return s.articleRepo.GetByID(ctx, id)
}

// CreateArticle creates a new article.
func (s *ArticleService) CreateArticle(ctx context.Context, article *models.Article) (*models.Article, error) {

	// initialize the upvotes and downvotes to zero
	article.UpVotes = 0
	article.DownVotes = 0

	// initialize the summary to the empty string
	article.Summary = ""

	// set the created at and updated at fields
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	return s.articleRepo.Create(ctx, article)
}

// DeleteArticle deletes an article by its ID.
func (s *ArticleService) DeleteArticle(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) error {

	// check if the article exists
	_, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return fmt.Errorf("%w : article not found", apperrors.ErrNotFound)
	}

	// check if the user is admin
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w : user not found", apperrors.ErrNotFound)
	}
	if !user.IsAdmin {
		return fmt.Errorf("%w : unauthorized to delete this article", apperrors.ErrUnauthorized)
	}

	// delete the comments associated with the article
	err = s.commentRepo.DeleteByArticleID(ctx, articleID)
	if err != nil {
		return fmt.Errorf("%w : failed to delete associated comments", apperrors.ErrInternal)
	}

	// delete the article
	return s.articleRepo.Delete(ctx, articleID)
}

func (s *ArticleService) SyncDailyArticles(ctx context.Context) error {
	// fetch the article from the client
	newArtciles, err := s.newsClient.FetchDailyArticles(ctx)
	if err != nil {
		return fmt.Errorf("%w : failed to fetch articles from news client", apperrors.ErrInternal)
	}

	errorCount := 0
	// save the articles in the database (maybe add CreateMany in the repository for better performance)
	for _, article := range newArtciles {
		_, err := s.articleRepo.Create(ctx, &article)
		if err != nil {
			log.Printf("Error saving article with title '%s' : %v\n", article.Title, err)
			errorCount++
			continue
		}
	}

	if errorCount > 0 {
		log.Printf("Failed to save %d articles for %d articles\n", errorCount, len(newArtciles))
	}

	return nil
}