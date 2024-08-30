package deploy

import (
	"gofr.dev/pkg/gofr"
	"kops.dev/models"
)

type client struct {
}

func New() *client {
	return &client{}
}

func (c *client) DeployImage(ctx *gofr.Context, img *models.Image) error {
	depSvc := ctx.GetHTTPService("deployment-service")

	resp, err := depSvc.Post(ctx, "/deploy", nil, getForm(img))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func getForm(img *models.Image) []byte {
	return nil
}
