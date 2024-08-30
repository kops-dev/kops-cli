package client

import (
	"context"

	"kops.dev/deployment-service/models"
)

type CredentialFetcher interface {
	GetServiceCreds(ctx context.Context, serviceID string) (*models.Credentials, error)
}

type ServiceImageUpdater interface {
	UpdateImage(ctx context.Context, serviceId, imageURL string) error
}
