package deploy

import (
	"fmt"
	"os"
	"path"

	"gofr.dev/pkg/gofr"
	"golang.org/x/sync/errgroup"

	"zop.dev/client"
	"zop.dev/models"
	"zop.dev/service"
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

	tempDirPath := path.Join(os.TempDir(), "zop-cli")

	err = ctx.File.Mkdir(tempDirPath, os.ModePerm)
	if err != nil {
		return err
	}

	defer os.RemoveAll(tempDirPath)

	zipFile, err := zipProject(img, tempDirPath)
	if err != nil {
		return err
	}

	err = s.depClient.Deploy(ctx, img, zipFile)
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
