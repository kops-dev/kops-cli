package build

import (
	"fmt"
	"os/exec"
)

const (
	golang = "golang"
	java   = "java"
	js     = "js"
)

// Need to build based on the language!
// Can we take the build command from the user?
// In case of node do we need to do the build using npx or npm or yarn?
// For every language support do we need to check if that language's compiler exists in the system?
// support - 1. golang(done)    2. Javascript      3. Java

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
		fmt.Println("error occurred while creating binary!", output)

		return err
	}

	fmt.Println("Binary created successfully")

	return nil
}
