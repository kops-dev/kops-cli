package deploy

import (
	"fmt"
	"os"

	"gofr.dev/pkg/gofr"
	"golang.org/x/sync/errgroup"

	"kops.dev/client"
	"kops.dev/models"
	"kops.dev/service"
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

func New(docker service.DockerClient, depClient client.ServiceDeployer) service.Deployer {
	return &svc{docker: docker, depClient: depClient}
}

func (s *svc) Deploy(ctx *gofr.Context, img *models.Image) error {
	err := buildProject(ctx)
	if err != nil {
		return err
	}

	// create the temp dir to save docker image that is built
	err = ctx.File.Mkdir("temp", os.ModePerm)
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

func buildProject(ctx *gofr.Context) error {
	lang := ctx.Param("lang")
	if lang == "" {
		lang = detect()
		if lang == "" {
			ctx.Logger.Errorf("%v", errLanguageNotProvided)

			return errLanguageNotProvided
		}

		fmt.Println("detected language is", lang)
	}

	group := errgroup.Group{}

	group.Go(func() error {
		err := Build(lang)
		if err != nil {
			ctx.Logger.Errorf("error while building the project binary, please check the project code!")

			return err
		}

		return nil
	})

	fi, _ := os.Stat("Dockerfile")
	if fi != nil {
		fmt.Println("Dockerfile present, using already created dockerfile")
	} else {
		group.Go(func() error {
			return createDockerFile(ctx, lang)
		})
	}

	return group.Wait()
}
