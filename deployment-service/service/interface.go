package service

import (
	"gofr.dev/pkg/gofr"

	"deployment-service/models"
)

type ImageUploader interface {
	UploadToArtifactory(ctx *gofr.Context, img *models.Image) (string, error)
}

type ImageDeployer interface {
	DeployImage(ctx *gofr.Context, serviceID, imageURL string) error
}
