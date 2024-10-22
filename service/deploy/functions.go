package deploy

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"gofr.dev/pkg/gofr"

	"zop.dev/models"
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
		fmt.Println("you can create your own DockerFile and run the zop.dev cli again.")

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

func zipProject(img *models.Image, zipDir string) (string, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	zipFile := path.Join(zipDir, fmt.Sprintf("%s-%s.zip", img.Name, img.Tag))

	outFile, err := os.Create(zipFile)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	err = filepath.Walk(curDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if file == zipFile {
			return nil
		}

		relPath, err := filepath.Rel(filepath.Dir(curDir), file)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			_, err = zipWriter.Create(relPath + "/")
			if err != nil {
				return err
			}

			return nil
		}

		fileInZip, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		_, err = io.Copy(fileInZip, fileToZip)

		return err
	})

	return zipFile, err
}
