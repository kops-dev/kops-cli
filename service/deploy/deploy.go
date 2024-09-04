package deploy

import (
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

type svc struct {
	docker    service.DockerClient
	depClient client.ServiceDeployer
}

func New(docker service.DockerClient, depClient client.ServiceDeployer) *svc {
	return &svc{docker: docker, depClient: depClient}
}

func (s *svc) Deploy(ctx *gofr.Context, img *models.Image) error {
	fi, _ := os.Stat("Dockerfile")
	if fi != nil {
		fmt.Println("Dockerfile present, using already created dockerfile")
	} else {
		if err := createDockerFile(ctx); err != nil {
			return err
		}
	}

	os.RemoveAll("temp")
	// create the temp dir to save docker image that is built
	err := ctx.File.Mkdir("temp", os.ModePerm)
	if err != nil {
		return err
	}

	defer os.RemoveAll("temp")

	err = s.docker.BuildImage(ctx, img)
	if err != nil {
		return err
	}

	err = s.docker.SaveImage(ctx, img)
	if err != nil {
		return err
	}

	err = zipImage(img)
	if err != nil {
		return err
	}

	err = s.depClient.DeployImage(ctx, img)
	if err != nil {
		return err
	}

	return nil
}
