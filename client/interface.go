package client

import (
	"gofr.dev/pkg/gofr"
	"zop.dev/models"
)

type ServiceDeployer interface {
	DeployImage(ctx *gofr.Context, img *models.Image) error
}
