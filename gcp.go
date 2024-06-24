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
	cmd := replaceInputOutput(
		exec.Command("gcloud", "auth", "activate-service-account", "--key-file=./application_creds.json", "--project="+key.ProjectId),
	)

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return err
	}

	defer func() {
		os.Remove("application_creds.json")
	}()

	cmd = replaceInputOutput(exec.Command("gcloud", "auth", "configure-docker", gcp.Region+"-docker.pkg.dev"))
	err = cmd.Run()
	if err != nil {
		fmt.Println("error configuring docker registry")
		return err
	}

	imageLoc := gcp.Region + "-docker.pkg.dev" + "/" + key.ProjectId + "/" + gcp.DockerRegistry + "/" + imageName

	cmd = replaceInputOutput(
		exec.Command("docker", "tag", imageName, imageLoc),
	)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = replaceInputOutput(
		exec.Command("docker", "push", imageLoc),
	)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("gcloud", "container", "clusters", "get-credentials", gcp.ClusterName, "--region="+gcp.Region, "--project="+key.ProjectId)
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("kubectl", "set", "image", "deployment/"+gcp.ServiceName,
		gcp.ServiceName+"="+imageLoc,
		"--namespace", gcp.Namespace)

	fmt.Println("kubectl", "set", "image", "deployment/"+gcp.ServiceName,
		gcp.ServiceName+"="+imageLoc,
		"--namespace", gcp.Namespace)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		fmt.Println("Detailed error output:", string(output))
		return err
	}

	return nil
}

// replaceInputOutput attaches the
func replaceInputOutput(cmd *exec.Cmd, files ...*os.File) *exec.Cmd {
	if len(files) > 0 {
		cmd.Stdout = files[0]
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd
}
