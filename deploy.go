package main

import (
	"errors"
	"fmt"
	"kops.dev/internal/templates"
	"os"
	"path/filepath"

	"gofr.dev/pkg/gofr"
)

var (
	errDepKeyNotProvided    = errors.New("KOPS_DEPLOYMENT_KEY not provided, please download the key form https://kops.dev")
	errLanguageNotProvided  = errors.New("unable to create DockerFile as project programming language not provided. Please Provide a programming language using -lang=<language>")
	errLanguageNotSupported = errors.New("creating DockerFile for provided language is not supported yet")
)

func Deploy(ctx *gofr.Context) (interface{}, error) {
	keyFile := os.Getenv("KOPS_DEPLOYMENT_KEY")
	if keyFile == "" {
		return nil, errDepKeyNotProvided
	}

	// letting this key file to be used later
	_, err := os.ReadFile(filepath.Clean(keyFile))
	if err != nil {
		return nil, err
	}

	fi, _ := os.Stat("Dockerfile")
	if fi != nil {
		fmt.Println("Dockerfile present, using already created dockerfile")
	} else {
		// removing the cloud-specific logic from cli to hosted service
		lang := ctx.Param("lang")
		if lang == "" {
			ctx.Logger.Errorf("%v", errLanguageNotProvided)

			return nil, errLanguageNotProvided
		}

		port := ctx.Param("p")
		if port == "" {
			port = "8000"
		}

		createDockerFile(ctx, lang)
	}

	return "Successful", nil
}

func createDockerFile(ctx *gofr.Context, lang string) error {
	content := templates.TmplMap[lang]
	if content == "" {
		ctx.Logger.Errorf("creating DockerFile for %s is not supported yet, reach us at https://github.com/kops-dev/kops-cli/issues to know more", lang)

		fmt.Printf("creating DockerFile for %s is not supported yet, reach us at https://github.com/kops-dev/kops-cli/issues to know more\n", lang)
		fmt.Println("you can create your own DockerFile and run the kops-cli again.")

		return errLanguageNotSupported
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		ctx.Logger.Error("error while creating DockerFile", err)
		fmt.Println("unable to create the DockerFile, please check permissions for creating files in the directory")

		return err
	}

	return nil
}
