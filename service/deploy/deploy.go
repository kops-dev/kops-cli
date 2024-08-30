package deploy

import (
	"errors"
	"fmt"
	"gofr.dev/pkg/gofr"
	"kops.dev/client"
	"kops.dev/models"
	"kops.dev/service"
	"os"
)

const (
	golang = "golang"
	java   = "java"
	js     = "js"
)

var (
	errDepKeyNotProvided = errors.New("KOPS_DEPLOYMENT_KEY not provided, " +
		"please download the key form https://kops.dev")
)

type svc struct {
	docker    service.DockerClient
	depClient client.ServiceDeployer
}

func New(docker service.DockerClient, depClient client.ServiceDeployer) *svc {
	return &svc{docker: docker, depClient: depClient}
}

func (s *svc) Deploy(ctx *gofr.Context, img *models.Image) error {
	// TODO: figure another approach with init or login commands
	//keyFile := os.Getenv("KOPS_DEPLOYMENT_KEY")
	//if keyFile == "" {
	//	return errDepKeyNotProvided
	//}
	//
	//_, err := os.ReadFile(filepath.Clean(keyFile))
	//if err != nil {
	//	return err
	//}

	fi, _ := os.Stat("Dockerfile")
	if fi != nil {
		fmt.Println("Dockerfile present, using already created dockerfile")
	} else {
		if err := createDockerFile(ctx); err != nil {
			return err
		}
	}

	err := s.docker.BuildImage(ctx, img)
	if err != nil {
		return err
	}

	err = s.docker.SaveImage(ctx, img)
	if err != nil {
		return err
	}

	err = s.depClient.DeployImage(ctx, img)
	if err != nil {
		return err
	}

	return nil
}
