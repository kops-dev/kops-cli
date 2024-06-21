package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"kops.dev/internal/models"
	"os"
	"os/exec"
)

var (
	errInvalidKey = errors.New("")
)

func deployGCP(gcp *models.Deploy, imageName string) error {
	var key models.GoogleCred

	jsonBytes, err := json.MarshalIndent(gcp.Key, "", " ")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())

		return errInvalidKey
	}

	// checks if the converted json data is valid
	err = json.Unmarshal(jsonBytes, &key)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())

		return errInvalidKey
	}

	err = os.WriteFile("application_creds.json", jsonBytes, 0644)
	if err != nil {
		return err
	}

	// Authenticate using gcloud
	cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file=./application_creds.json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))

		return err
	}

	fmt.Println(string(out))

	defer func() {
		fmt.Println("removing application credentials")
		os.Remove("application_creds.json")
	}()

	return nil
}
