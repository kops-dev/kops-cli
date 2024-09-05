package deploy

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"gofr.dev/pkg/gofr"

	"kops.dev/models"
)

const (
	imageZipName = "temp/image.zip"
)

var (
	errLanguageNotProvided = errors.New("unable to create DockerFile as project " +
		"programming language not provided. Please Provide a programming language using -lang=<language>")
	errLanguageNotSupported = errors.New("creating DockerFile for provided language is not supported yet")
)

func createDockerFile(ctx *gofr.Context, lang string) error {
	var content, port string

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

// TODO: For every language support do we need to check if that language's compiler exists in the system.
// support - 1. golang(done)    2. Javascript      3. Java

// Build executes the build command for the project specific to language.
func Build(lang string) error {
	switch lang {
	case golang:
		return buildGolang()
	case js:
		// TODO: necessary steps for javascript build
		break
	case java:
		// TODO: necessary steps for building java projects
		break
	}

	return nil
}

func buildGolang() error {
	fmt.Println("Creating binary for the project")

	output, err := exec.Command("sh", "-c", "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .").CombinedOutput()
	if err != nil {
		fmt.Println("error occurred while creating binary!", string(output))

		return err
	}

	fmt.Println("Binary created successfully")

	return nil
}
