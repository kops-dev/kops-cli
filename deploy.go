package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gofr.dev/pkg/gofr"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"kops.dev/internal/models"
)

const (
	gcp            = "GCP"
	executableName = "main"
)

var (
	buildFlags = fmt.Sprintf("CGO_ENABLED=0 GOOS=%s GOARCH=%s", runtime.GOOS, runtime.GOARCH)
)

func Deploy(ctx *gofr.Context) (interface{}, error) {
	var deploy models.Deploy

	keyFile := os.Getenv("KOPS_DEPLOYMENT_KEY")
	if keyFile == "" {
		return nil, errors.New("KOPS_DEPLOYMENT_KEY not provided, please download the key form https://kops.dev")
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

	image := deploy.ServiceName + ":" + tag
	cmd := exec.Command("docker", "build", "-t"+image, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// check what cloud provider is
	switch deploy.CloudProvider {
	default:
		err = deployGCP(&deploy, image)
		if err != nil {
			return nil, err
		}

		return "Successfully deployed!", nil
	}
}

func buildBinary() error {
	cmd := exec.Command("sh", "-c", buildFlags+" go build -o "+executableName+" .")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing command: %s, %v", output, err)
	}

	fmt.Println("Binary created successfully")

	return nil
}

func createGoDockerFile() error {
	content := fmt.Sprintf(`FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
COPY main ./main
RUN chmod +x /main
EXPOSE 8000
CMD ["/main"]`)

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
