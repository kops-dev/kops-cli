package upload

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"

	"kops.dev/deployment-service/client"
	"kops.dev/deployment-service/models"
)

type service struct {
	credSvc client.CredentialFetcher
}

func New(credSvc client.CredentialFetcher) *service {
	return &service{credSvc: credSvc}
}

func (s *service) UploadToArtifactory(ctx *gofr.Context, img *models.Image) (string, error) {
	dir := getUniqueDir()

	err := img.Data.CreateLocalCopies(dir)
	if err != nil {
		return "", err
	}

	creds, err := s.credSvc.GetServiceCreds(ctx, img.ServiceID)
	if err != nil {
		return "", err
	}

	path, err := pushImage(ctx, img, creds, dir)
	if err != nil {
		return "", err
	}

	ctx.Logger.Infof("successfully pushed image %v to artifact registry", img.Name)

	return path, nil
}

func pushImage(ctx *gofr.Context, img *models.Image, cred *models.Credentials, path string) (string, error) {
	var (
		imagePath string
		auth      *authn.Basic
	)

	switch cred.CloudPlatform {
	case GCP:
		googleReg, err := NewGCR(cred)
		if err != nil {
			return "", err
		}

		imagePath = googleReg.getImagePath(img)
		auth, err = googleReg.getAuth(ctx)
		if err != nil {
			return "", err
		}
	}

	ref, err := name.ParseReference(imagePath)
	if err != nil {
		return "", err
	}

	imgTar, err := tarball.ImageFromPath(path, nil)
	if err != nil {
		return "", err
	}

	// Push the image to the specified registry
	err = remote.Write(ref, imgTar, remote.WithAuth(auth))
	if err != nil {
		return "", err
	}

	return imagePath, nil
}

func getUniqueDir() string {
	dirName, _ := uuid.NewUUID()
	return dirName.String()
}
