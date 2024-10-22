package service

import (
	"gofr.dev/pkg/gofr"
	"zop.dev/models"
)

type Deployer interface {
	Deploy(ctx *gofr.Context, img *models.Image) error
}
