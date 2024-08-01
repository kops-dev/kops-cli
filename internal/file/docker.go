package file

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"gofr.dev/pkg/gofr"

	"kops.dev/internal/templates"
)

var errLanguageNotSupported = errors.New("creating DockerFile for provided language is not supported yet")

func CreateDockerFile(ctx *gofr.Context, lang, port string) error {
	var content string

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
