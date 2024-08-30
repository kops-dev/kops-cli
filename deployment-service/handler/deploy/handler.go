package deploy

import (
	"gofr.dev/pkg/gofr"

	"kops.dev/deployment-service/models"
	"kops.dev/deployment-service/service"
)

type handler struct {
	uploadSvc service.ImageUploader
	deploySvc service.ImageDeployer
}

func New(uploadSvc service.ImageUploader, deploySvc service.ImageDeployer) *handler {
	return &handler{
		uploadSvc: uploadSvc,
		deploySvc: deploySvc,
	}
}

func (h *handler) UploadImage(ctx *gofr.Context) (interface{}, error) {
	var img models.Image

	if err := ctx.Bind(&img); err != nil {
		return nil, err
	}

	// call the uploader service to
	imageURL, err := h.uploadSvc.UploadToArtifactory(ctx, &img)
	if err != nil {
		return nil, err
	}

	err = h.deploySvc.DeployImage(ctx, img.ServiceID, imageURL)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
