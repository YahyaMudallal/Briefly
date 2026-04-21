package clients

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type GeminiClient interface {
	GenerateTLDR(ctx context.Context, articleID bson.ObjectID) (string, error)
}

type GeminiAPIClient struct {
	apiKey string
}

func NewGeminiAPIClient(apiKey string) GeminiClient {
	return &GeminiAPIClient{apiKey: apiKey}
}

func (c *GeminiAPIClient) GenerateTLDR(ctx context.Context, articleID bson.ObjectID) (string, error) {
	// TODO
	return "", nil
}