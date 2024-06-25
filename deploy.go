package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gofr.dev/pkg/gofr"
	"kops.dev/internal/models"
)

const (
	gcp            = "GCP"
	executableName = "main"
)

var (
	errDepKeyNotProvided          = errors.New("KOPS_DEPLOYMENT_KEY not provided, please download the key form https://kops.dev")
	errCloudProviderNotRecognized = errors.New("cloud provider in KOPS_DEPLOYMENT_KEY is not provided or supported")
)

func Deploy(ctx *gofr.Context) (interface{}, error) {
	var deploy models.Deploy

	keyFile := os.Getenv("KOPS_DEPLOYMENT_KEY")
	if keyFile == "" {
		return nil, errDepKeyNotProvided
	}

	f, err := os.ReadFile(filepath.Clean(keyFile))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(f, &deploy)
	if err != nil {
		return nil, err
	}

	// build binaries for the current working directory
	err = buildBinary()
	if err != nil {
		return nil, err
	}

	// create and build docker image
	err = createGoDockerFile()
	if err != nil {
		return nil, err
	}

	tag := ctx.Param("tag")
	if tag == "" {
		tag = "latest"
	}

	image := filepath.Clean(deploy.ServiceName + ":" + tag)

	err = replaceInputOutput(exec.Command("docker", "build", "-t", image, ".")).Run()
	if err != nil {
		return nil, err
	}

	// check what cloud provider is
	switch deploy.CloudProvider {
	case gcp:
		err = deployGCP(&deploy, image)
		if err != nil {
			return nil, err
		}

		return "Successfully deployed!", nil
	default:
		return nil, errCloudProviderNotRecognized
	}
}

func buildBinary() error {
	fmt.Println("Creating binary for the project")

	output, err := exec.Command("sh", "-c", "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "+executableName+" .").CombinedOutput()
	if err != nil {
		fmt.Println("error occurred while creating binary!", output)

		return err
	}

	fmt.Println("Binary created successfully")

	return nil
}

func createGoDockerFile() error {
	content := `FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
COPY main ./main
RUN chmod +x /main
EXPOSE 8000
CMD ["/main"]`

	fi, _ := os.Stat("Dockerfile")
	if fi != nil {
		fmt.Println("Dockerfile present, using already created dockerfile")
		return nil
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		return err
	}

	fmt.Println("Dockerfile created successfully!")

	return nil
}
