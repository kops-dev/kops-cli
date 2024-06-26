package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"kops.dev/internal/models"
)

const filePerm = 0644

var (
	errInvalidKey = errors.New("key provided is not a valid service-account key")
)

func deployGCP(gcp *models.Deploy, imageName string) error {
	var key = models.GCPInfo{}

	jsonBytes, err := json.MarshalIndent(gcp.Key, "", " ")
	if err != nil {
		fmt.Println(err.Error())

		return errInvalidKey
	}

	err = authenticateCLI(jsonBytes, &key)
	if err != nil {
		return err
	}

	//nolint:gosec // this would not be an issue since the user input is part of secret
	err = replaceInputOutput(exec.Command("gcloud", "auth",
		"configure-docker", gcp.Region+"-docker.pkg.dev")).Run()
	if err != nil {
		fmt.Println("error configuring docker registry")
		return err
	}

	imageLoc := gcp.Region + "-docker.pkg.dev" + "/" + key.ProjectID + "/" + gcp.DockerRegistry + "/" + imageName

	err = replaceInputOutput(
		exec.Command("docker", "tag", imageName, imageLoc),
	).Run()
	if err != nil {
		return err
	}

	err = replaceInputOutput(
		exec.Command("docker", "push", imageLoc),
	).Run()
	if err != nil {
		return err
	}

	//nolint:gosec // this would not be an issue since the user input is part of secret
	err = replaceInputOutput(exec.Command("gcloud", "container", "clusters",
		"get-credentials", gcp.ClusterName, "--region="+gcp.Region, "--project="+key.ProjectID)).Run()
	if err != nil {
		return err
	}

	//nolint:gosec // this would not be an issue since the user input is part of secret
	err = replaceInputOutput(exec.Command("kubectl", "set", "image", "deployment/"+gcp.ServiceName,
		gcp.ServiceName+"="+imageLoc,
		"--namespace", gcp.Namespace)).Run()
	if err != nil {
		return err
	}

	return nil
}

// replaceInputOutput attaches the current terminal std out, err and in to cmd.
func replaceInputOutput(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd
}

func authenticateCLI(jsonBytes []byte, info *models.GCPInfo) error {
	var key models.GCPInfo

	// checks if the converted json data is valid
	err := json.Unmarshal(jsonBytes, &key)
	if err != nil {
		fmt.Println(err.Error())

		return errInvalidKey
	}

	info.ProjectID = key.ProjectID

	f, err := os.OpenFile("application_creds.json", os.O_CREATE|os.O_WRONLY, filePerm)
	if err != nil {
		return err
	}

	defer os.Remove("application_creds.json")

	_, err = f.Write(jsonBytes)
	if err != nil {
		return err
	}

	// authenticate gcloud cli with the newly created credentials, these are part of the KOPS_DEPLOYMENT_KEY
	//nolint:gosec // this would not be an issue since the user input is part of secret
	err = replaceInputOutput(
		exec.Command("gcloud", "auth",
			"activate-service-account", "--key-file=./application_creds.json", "--project="+key.ProjectID)).Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return err
	}

	return nil
}
