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
	userRepo    repositories.UserRepository
	commentRepo repositories.CommentRepository
	newsClient clients.NewsClient
	geminiClient clients.GeminiClient
  votesRepo   repositories.VoteRepository
}

// NewArticleService creates a new ArticleService.
func NewArticleService(
	articleRepo repositories.ArticleRepository,
	userRepo repositories.UserRepository,
	commentRepo repositories.CommentRepository,
	newsClient clients.NewsClient,
	geminiClient clients.GeminiClient,
  	votesRepo repositories.VoteRepository,
	) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		userRepo: userRepo,
		commentRepo: commentRepo,
		newsClient: newsClient,
		geminiClient: geminiClient,
    	votesRepo: votesRepo,
	}
}

// Create a new Response Struct that combines them
type ArticleResponse struct {
	models.Article `bson:",inline"` // Embeds all the standard article fields
	UserVote       int              `json:"userVote"`
}

// GetPaginated retrieves a paginated list of articles.
func (s *ArticleService) GetPaginated(
	ctx context.Context,
	page int,
	limit int,
	sortBy string,
	order string) ([]ArticleResponse, error) {
	//we need to inject logged-in user vote status into the articles before returning them,
	//so we need to get all the articles first, then get the votes for each article for the logged-in user and inject them into the articles

	articles, err := s.articleRepo.GetPaginated(ctx, page, limit, sortBy, order)
	if err != nil {
		return nil, fmt.Errorf("%w : failed to retrieve articles", apperrors.ErrInternal)
	}
	userID, ok := ctx.Value("user_id").(bson.ObjectID)

	var userVotes map[bson.ObjectID]int // Map to easily look up votes by Article ID

	if userID != bson.NilObjectID && ok {
		// Get all votes for the logged-in user for these articles
		votes, err := s.votesRepo.GetAllByUserID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("%w : failed to retrieve user votes", apperrors.ErrInternal)
		}

		// Build the map
		userVotes = make(map[bson.ObjectID]int)
		for _, v := range votes {
			userVotes[v.ArticleID] = v.Type // 1 or -1
		}
	}

	var feed []ArticleResponse
	for _, article := range articles {
		feed = append(feed, ArticleResponse{
			Article:  article,
			UserVote: userVotes[article.ID], // Will be 0 if the user hasn't voted on this article
		})
	}

	return feed, nil
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

// SyncDailyArticles Called daily by the server to Synchronize the articles database
// by calling the new api to get new articles and add them to the database.
func (s *ArticleService) SyncDailyArticles(ctx context.Context) error {
	// fetch the article from the client
	newArticles, err := s.newsClient.FetchDailyArticles(ctx)
	if err != nil {
		return fmt.Errorf("%w : failed to fetch articles from news client", apperrors.ErrInternal)
	}

	log.Printf("DEBUG: The news API returned %d articles", len(newArticles))

	// save the articles in the database (maybe add CreateMany in the repository for better performance)
	_, err = s.articleRepo.CreateMany(ctx, newArticles)
	if err != nil {
		return fmt.Errorf("%w : failed to save articles in the database", apperrors.ErrInternal)
	}

	return nil
}

// GenerateSummary call the Gemini API to generate a summary (tldr) for the given article and update the article in the database with the new summary.
func (s *ArticleService) GenerateSummary(ctx context.Context, articleID bson.ObjectID) (string, error) {

	// get the article from the database
	article, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return "", fmt.Errorf("%w : failed to get the article from the database: %w", apperrors.ErrInternal, err)
	}

	// check if the article has already a summary
	if article.Summary != "" {
		return "", fmt.Errorf("%w : the article already has a summary", apperrors.ErrValidation)
	}

	// call the Gemini API to generate a summary for the article
	summary, err := s.geminiClient.GenerateTLDR(ctx, article)
	if err != nil {
		return "", fmt.Errorf("%w : failed to generate the summary: %w", apperrors.ErrInternal, err)
	}

	// add the summary to the article
	article.Summary = summary

	// update the updated at date
	article.UpdatedAt = time.Now()

	// update the article in the database
	err = s.articleRepo.Update(ctx, article)
	if err != nil {
		return "", fmt.Errorf("%w : failed to update the article in the database: %w", apperrors.ErrInternal, err)
	}

	return summary, nil
}

// ToggleUpvote toggles the upvote status of an article by its ID.
func (s *ArticleService) ToggleUpvote(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) (int, error) {
	// check if the article exists
	_, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return 0, fmt.Errorf("%w : article not found", apperrors.ErrNotFound)
	}

	// ToggleUpvote returns 1 if it was from an upvote to no vote, -1 if it was from a downvote to upvote, and 0 if it was from no vote to upvote
	prevStatus, err := s.votesRepo.ToggleUpvote(ctx, articleID, userID)
	if err != nil {
		return 0, fmt.Errorf("%w : failed to toggle upvote", apperrors.ErrInternal)
	}
	//update the article's upvote count
	switch prevStatus {
	case 0: // novote -> upvote
		err = s.articleRepo.IncrementUpVotes(ctx, articleID, 1)
	case -1: // downvote -> upvote
		err = s.articleRepo.IncrementUpVotes(ctx, articleID, 1)
		err = s.articleRepo.IncrementDownVotes(ctx, articleID, -1)
	case 1: // upvote -> novote
		err = s.articleRepo.IncrementUpVotes(ctx, articleID, -1)
	}
	if err != nil {
		return 0, fmt.Errorf("%w : failed to increment article upvote count", apperrors.ErrInternal)
	}
	article, err := s.articleRepo.GetByID(ctx, articleID)

	return article.UpVotes, nil
}

// ToggleDownvote toggles the downvote status of an article by its ID.
func (s *ArticleService) ToggleDownvote(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) (int, error) {
	// check if the article exists
	_, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return 0, fmt.Errorf("%w : article not found", apperrors.ErrNotFound)
	}

	prevStatus, err := s.votesRepo.ToggleDownvote(ctx, articleID, userID)
	if err != nil {
		return 0, fmt.Errorf("%w : failed to toggle downvote", apperrors.ErrInternal)
	}
	// update the article's downvote count
	switch prevStatus {
	case 0: // novote -> downvote
		err = s.articleRepo.IncrementDownVotes(ctx, articleID, 1)
	case 1: // upvote -> downvote
		err = s.articleRepo.IncrementDownVotes(ctx, articleID, 1)
		err = s.articleRepo.IncrementUpVotes(ctx, articleID, -1)
	case -1: // downvote -> novote
		err = s.articleRepo.IncrementDownVotes(ctx, articleID, -1)
	}
	if err != nil {
		return 0, fmt.Errorf("%w : failed to increment article downvote count", apperrors.ErrInternal)
	}
	article, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return 0, fmt.Errorf("%w : failed to get article", apperrors.ErrInternal)
	}
	return article.DownVotes, nil
}
