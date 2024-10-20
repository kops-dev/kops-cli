package client

import (
	"gofr.dev/pkg/gofr"
	"zop.dev/models"
)

type ServiceDeployer interface {
	Deploy(ctx *gofr.Context, img *models.Image, zipFile string) error
}
