package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YahyaMudallal/newsWebSite/internal/models"
)

// base url for the news API
const newAPIEndpoint = "https://newsdata.io/api/1/latest"

type NewsClient interface {
    FetchDailyArticles(ctx context.Context) ([]models.Article, error)
}	

type newsAPIClient struct {
    apiKey string
	httpClient *http.Client
}

func NewNewsAPIClient(apiKey string) NewsClient {
    return &newsAPIClient{
		apiKey: apiKey,
		httpClient: &http.Client{Timeout: 15 * time.Second}}	// 15 seconds timeout for security
}

type newsAPIResponse struct {
	Status			string				`json:"status"`
	TotalResults 	int					`json:"totalResults"`
	Results 		[]newsAPIArticle 	`json:"results"`
}

// newsAPIArticle represents the structure of an article returned by the API
type newsAPIArticle struct {
	Title			string `json:"title"`
	Description 	string `json:"description"`
	Content			string `json:"content"`
	Link			string `json:"link"`
	ImageURL		string `json:"image_url"`
	PubDate		 	string `json:"pubDate"`
	SourceID 		string `json:"source_id"`
}

// FetchDailyArticles fetches the daily articles from the News API.
func (c *newsAPIClient) FetchDailyArticles(ctx context.Context) ([]models.Article, error) {
	
	// create a new request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, newAPIEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request : %w", err)
	}

	// add the query parameters
	q := req.URL.Query()
	q.Add("apikey", c.apiKey)
	q.Add("language", "en")
	q.Add("size", "10")
	//q.Add("image", "1") 	// uncomment to force only article with images
	req.URL.RawQuery = q.Encode()

	// make the HTTP request
	resp, err := c.httpClient.Do(req)	
	if err != nil {
		return nil, fmt.Errorf("failed to make request : %w", err)
	}
	defer resp.Body.Close()

	// check http status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("news API returned status : %d", resp.StatusCode)
	}

	// decode the response
	var apiResp newsAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the response : %w", err)
	}

	// convert struct of api to struct of model
	var articles []models.Article
	for _, item := range apiResp.Results {

		publishedDate, err := time.Parse("2006-01-02 15:04:05", item.PubDate)
		if err != nil {
			publishedDate = time.Now()		// security fallback
		}

		// create the article model
		article := models.Article{
			Title: item.Title,
			Description: item.Description,
			Content: item.Content,
			OriginalURL: item.Link,
			ImageURL: item.ImageURL,
			Source: item.SourceID,
			PublishedAt: publishedDate,
		}

		// fallback if content is empty
        if article.Content == "" {
            article.Content = article.Description
        }

		// the article to the list
		articles = append(articles, article)
	}

	return articles, nil
}