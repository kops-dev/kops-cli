package deploy

import (
	"gofr.dev/pkg/gofr"

	"kops.dev/models"
	"kops.dev/service"
)

type handler struct {
	svc service.Deployer
}

func New(svc service.Deployer) *handler {
	return &handler{svc: svc}
}

func (h *handler) Deploy(ctx *gofr.Context) (any, error) {
	img := &models.Image{
		Name:       ctx.Param("name"),
		Tag:        ctx.Param("tag"),
		ServiceID:  ctx.Param("service"),
		Repository: ctx.Param("repository"),
		Region:     ctx.Param("region"),
	}

	err := h.svc.Deploy(ctx, img)
	if err != nil {
		return nil, err
	}

	return "Successfully deployed " + img.Name, nil
}
