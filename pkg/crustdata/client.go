package crustdata

import (
	"context"
)

// CrustdataServiceInterface defines the contract for Crustdata API operations
type CrustdataServiceInterface interface {
	// People Discovery API
	PeopleSearch(ctx context.Context, req *PeopleSearchRequest) (*PeopleSearchResponse, error)
}

// CrustdataService implements the CrustdataServiceInterface
type CrustdataService struct {
	Client  *CrustdataClient
	BaseURL string
	APIKey  string
}

// NewCrustdataService creates a new Crustdata service instance
func NewCrustdataService(apiKey string, baseURL string) (CrustdataServiceInterface, error) {
	client := NewClient(apiKey, baseURL)
	return &CrustdataService{
		Client:  client,
		BaseURL: baseURL,
		APIKey:  apiKey,
	}, nil
}
