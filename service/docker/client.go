package docker

import (
	"encoding/json"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"gofr.dev/pkg/gofr"

	"kops.dev/models"
	kopsSvc "kops.dev/service"
)

type service struct {
	docker *client.Client
}

func New() kopsSvc.DockerClient {
	c, err := client.NewClientWithOpts(
		client.WithHost("unix:///Users/raramuri/.colima/default/docker.sock"),
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return nil
	}

	return &service{docker: c}
}

type buildOutput struct {
	Stream string `json:"stream"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func (s *service) BuildImage(ctx *gofr.Context, img *models.Image) error {
	buildContext, err := archive.TarWithOptions(".", &archive.TarOptions{})
	if err != nil {
		ctx.Errorf("unable to generate the build context for current project, error : %v", err)

		return err
	}

	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		Tags:           []string{img.Name + ":" + img.Tag},
		Dockerfile:     "/Dockerfile",
	}

	imageBuildResponse, err := s.docker.ImageBuild(ctx, buildContext, options)
	if err != nil {
		ctx.Errorf("unable to build docker image, error : %v", err)

		return err
	}

	defer imageBuildResponse.Body.Close()

	// Decode and print formatted output
	decoder := json.NewDecoder(imageBuildResponse.Body)

	for {
		var output buildOutput
		if er := decoder.Decode(&output); er == io.EOF {
			break
		} else if er != nil {
			ctx.Error(er)
		}

		if output.Stream != "" && output.Stream != `\n'` {
			ctx.Debug(output.Stream)
		}

		if output.Status != "" {
			ctx.Info(output.Status)
		}

		if output.Error != "" {
			ctx.Debug("Error: %s\n", output.Error)
		}
	}

	return nil
}

func (s *service) SaveImage(ctx *gofr.Context, img *models.Image) error {
	tarFileName := "temp/" + img.Name + img.Tag + ".tar"

	tarFile, err := os.Create(tarFileName)
	if err != nil {
		ctx.Errorf("unable to create the image tar file, error : %v", err)
		return err
	}
	defer tarFile.Close()

	reader, err := s.docker.ImageSave(ctx, []string{img.Name + ":" + img.Tag})
	if err != nil {
		ctx.Errorf("uanble to create save image, error : %v", err)
		return err
	}
	defer reader.Close()

	// Write the image data to the tar file
	_, err = io.Copy(tarFile, reader)
	if err != nil {
		ctx.Logger.Errorf("unable to save image to file, error : %v", err)
		return err
	}

	return nil
}
