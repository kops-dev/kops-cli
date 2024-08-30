package deploy

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"kops.dev/models"
	"os"
	"strings"
	"text/template"

	"gofr.dev/pkg/gofr"
)

const (
	imageZipName = "image.zip"
)

var (
	errLanguageNotProvided = errors.New("unable to create DockerFile as project " +
		"programming language not provided. Please Provide a programming language using -lang=<language>")
	errLanguageNotSupported = errors.New("creating DockerFile for provided language is not supported yet")
)

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
		content = Golang
	case java:
		content = Java
	case js:
		content = Js
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

func zipImage(img *models.Image) error {
	iamgeTarName := "temp/" + img.Name + img.Tag + ".tar"

	tarReader, err := os.Open(iamgeTarName)
	if err != nil {
		return err
	}
	defer tarReader.Close()

	// Create the zip file for writing
	zipWriter, err := os.Create(imageZipName)
	if err != nil {
		return err
	}
	defer zipWriter.Close()

	archive := zip.NewWriter(zipWriter)
	defer archive.Close()

	w, err := archive.Create(iamgeTarName)
	if err != nil {
		return err
	}

	// Copy the tar file content to the zip writer
	_, err = io.Copy(w, tarReader)
	if err != nil {
		return err
	}

	return nil
}
