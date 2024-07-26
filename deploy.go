package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gofr.dev/pkg/gofr"

	"kops.dev/internal/templates"
)

const (
	golang = "golang"
	java   = "java"
	js     = "js"
)

var (
	errDepKeyNotProvided = errors.New("KOPS_DEPLOYMENT_KEY not provided, " +
		"please download the key form https://kops.dev")
	errLanguageNotProvided = errors.New("unable to create DockerFile as project " +
		"programming language not provided. Please Provide a programming language using -lang=<language>")
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
		if err := createDockerFile(ctx); err != nil {
			return nil, err
		}
	}

	// TODO: build and push the docker image to the Kops API
	// Also need to figure out the contract for the API

	return "Successful", nil
}

func createDockerFile(ctx *gofr.Context) error {
	var content, lang, port string

	// removing the cloud-specific logic from cli to hosted service
	lang = ctx.Param("lang")
	if lang == "" {
		lang = detect()
		if lang == "" {
			ctx.Logger.Errorf("%v", errLanguageNotProvided)

			return errLanguageNotProvided
		}

		fmt.Println("detected language is", lang)
	}

	port = ctx.Param("p")
	if port == "" {
		port = "8000"
	}

	// get the template content for dockerFile based on the language
	switch strings.ToLower(lang) {
	case golang:
		content = templates.Golang
	case java:
		content = templates.Java
	case js:
		content = templates.Js
	default:
		ctx.Logger.Errorf("creating DockerFile for %s is not supported yet,"+
			" reach us at https://github.com/kops-dev/kops-cli/issues to know more", lang)

		fmt.Printf("creating DockerFile for %s is not supported yet, "+
			"reach us at https://github.com/kops-dev/kops-cli/issues to know more\n", lang)
		fmt.Println("you can create your own DockerFile and run the kops-cli again.")

		return errLanguageNotSupported
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	defer file.Close()

	t := template.New(lang)

	temp, err := t.Parse(content)
	if err != nil {
		return err
	}

	if er := temp.Execute(file, port); er != nil {
		ctx.Logger.Error("error while creating DockerFile", er)
		fmt.Println("unable to create the DockerFile, please check permissions for creating files in the directory")

		return er
	}

	return nil
}

func detect() string {
	switch {
	case checkFile("go.mod"):
		return golang
	case checkFile("package.json"):
		return js
	case checkFile("pom.xml") || checkFile("build.gradle"):
		return java
	}

	return ""
}

func checkFile(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return false
	}

	return true
}
