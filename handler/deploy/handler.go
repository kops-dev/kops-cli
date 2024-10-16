package deploy

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gofr.dev/pkg/gofr"

	"zop.dev/models"
	"zop.dev/service"
)

var (
	errDepKeyNotProvided = errors.New("KOPS_DEPLOYMENT_KEY not provided, " +
		"please download the key form https://zop.dev and navigating to your service deployment guide lines for CLI deployment")
	errIncorrectDepKey = errors.New("unable to validate the deployment key, please make sure the key contents of key are correct" +
		"Please download the key form https://zop.dev and navigating to your service deployment guide lines for CLI deployment")
	errDeploymentFailed = errors.New("some unexpected error occurred while deploying your service using zop.dev")
)

type Handler struct {
	svc service.Deployer
}

func New(svc service.Deployer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Deploy(ctx *gofr.Context) (any, error) {
	var img models.Image

	keyFile := os.Getenv("ZOP_DEPLOYMENT_KEY")
	if keyFile == "" {
		return nil, errDepKeyNotProvided
	}

	f, err := os.ReadFile(filepath.Clean(keyFile))
	if err != nil {
		ctx.Logger.Errorf("error reading the deployment key")

		return nil, err
	}

	err = json.Unmarshal(f, &img)
	if err != nil {
		ctx.Logger.Errorf("error while binding the image details, err : %v", err)

		return nil, errIncorrectDepKey
	}

	img.Name = ctx.Param("name")
	img.Tag = ctx.Param("tag")

	err = h.svc.Deploy(ctx, &img)
	if err != nil {
		ctx.Logger.Errorf("error updating the user service, err : %v", err)

		return nil, errDeploymentFailed
	}

	return "Successfully deployed " + img.Name, nil
}
