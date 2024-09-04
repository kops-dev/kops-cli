package deploy

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gofr.dev/pkg/gofr"

	"kops.dev/models"
	"kops.dev/service"
)

var (
	errDepKeyNotProvided = errors.New("KOPS_DEPLOYMENT_KEY not provided, " +
		"please download the key form https://kops.dev and navigating to your service deployment guide lines for CLI deployment")
)

type Handler struct {
	svc service.Deployer
}

func New(svc service.Deployer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Deploy(ctx *gofr.Context) (any, error) {
	var img models.Image

	keyFile := os.Getenv("KOPS_DEPLOYMENT_KEY")
	if keyFile == "" {
		return nil, errDepKeyNotProvided
	}

	f, err := os.ReadFile(filepath.Clean(keyFile))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(f, &img)
	if err != nil {
		return nil, err
	}

	img.Name = ctx.Param("name")
	img.Tag = ctx.Param("tag")

	err = h.svc.Deploy(ctx, &img)
	if err != nil {
		return nil, err
	}

	return "Successfully deployed " + img.Name, nil
}
