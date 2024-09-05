package service

import (
	"gofr.dev/pkg/gofr"
	"kops.dev/models"
)

type Deployer interface {
	Deploy(ctx *gofr.Context, img *models.Image) error
}

type DockerClient interface {
	BuildImage(ctx *gofr.Context, img *models.Image) error
	SaveImage(ctx *gofr.Context, img *models.Image) error
}
