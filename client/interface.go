package client

import (
	"gofr.dev/pkg/gofr"
	"kops.dev/models"
)

type ServiceDeployer interface {
	DeployImage(ctx *gofr.Context, img *models.Image) error
}
