package upload

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"gofr.dev/pkg/gofr"
	"golang.org/x/oauth2/google"

	"kops.dev/deployment-service/models"
)

const GCP = "GCP"

var errIncorrectCloud = errors.New("")

type gcr struct {
	creds *models.GCPCreds
}

func NewGCR(creds *models.Credentials) (*gcr, error) {
	gcpCreds, ok := creds.ServiceAccCred.(models.GCPCreds)
	if !ok {
		return nil, errIncorrectCloud
	}

	return &gcr{creds: &gcpCreds}, nil
}

func (g *gcr) getImagePath(img *models.Image) string {
	return fmt.Sprintf("%s-docker.pkg.dev/%s/%s/%s:%s", img.Region, g.creds.ProjectID, img.Repository, img.Name, img.Tag)
}

func (g *gcr) getAuth(ctx *gofr.Context) (*authn.Basic, error) {
	b, _ := json.Marshal(g.creds)

	creds, err := google.CredentialsFromJSON(ctx, b, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	return &authn.Basic{
		Username: "oauth2accesstoken",
		Password: token.AccessToken,
	}, nil
}
