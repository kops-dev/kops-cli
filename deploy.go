package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gofr.dev/pkg/gofr"

	"kops.dev/internal/build"
	"kops.dev/internal/file"
)

var (
	errLanguageNotProvided = errors.New("unable to create DockerFile as project " +
		"programming language not provided. Please Provide a programming language using -lang=<language>")
	errDepKeyNotProvided = errors.New("KOPS_DEPLOYMENT_KEY not provided, " +
		"please download the key form https://kops.dev")
)

func Deploy(ctx *gofr.Context) (interface{}, error) {
	var lang, port string

	keyFile := os.Getenv("KOPS_DEPLOYMENT_KEY")
	if keyFile == "" {
		return nil, errDepKeyNotProvided
	}

	// letting this key file to be used later
	_, err := os.ReadFile(filepath.Clean(keyFile))
	if err != nil {
		return nil, err
	}

	// removing the cloud-specific logic from cli to hosted service
	lang = ctx.Param("lang")
	if lang == "" {
		lang = file.Detect()
		if lang == "" {
			ctx.Logger.Errorf("%v", errLanguageNotProvided)

			return nil, errLanguageNotProvided
		}

		fmt.Println("Detected language is", lang)
	}

	port = ctx.Param("p")
	if port == "" {
		port = "8000"
	}

	fi, _ := os.Stat("Dockerfile")
	if fi != nil {
		fmt.Println("Dockerfile present, using already created dockerfile")
	} else {
		if er := file.CreateDockerFile(ctx, lang, port); er != nil {
			return nil, er
		}
	}

	// Build the binary based on the language detected or provided by user
	err = build.Build(lang)
	if err != nil {
		ctx.Logger.Errorf("error building for %v, error : %v", lang, err)

		return nil, err
	}

	// TODO: Need to figure out the contract for the API

	return "Successful", nil
}
