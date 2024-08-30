package deploy

import (
	"gofr.dev/pkg/gofr"

	"kops.dev/deployment-service/client"
)

type service struct {
	svcImageUpdater client.ServiceImageUpdater
}

func New(siu client.ServiceImageUpdater) *service {
	return &service{svcImageUpdater: siu}
}

func (s *service) DeployImage(ctx *gofr.Context, serviceID, imageURL string) error {
	// TODO: validate serviceID and fetch its details
	// TODO: if service is not present then we can create a new service and then deploy

	err := s.svcImageUpdater.UpdateImage(ctx, serviceID, imageURL)
	if err != nil {
		return err
	}

	return nil
}
