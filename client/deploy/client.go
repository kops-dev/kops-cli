package deploy

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"

	"gofr.dev/pkg/gofr"
	"kops.dev/models"
)

const (
	imageZipName = "temp/image.zip"
)

type client struct {
}

func New() *client {
	return &client{}
}

func (c *client) DeployImage(ctx *gofr.Context, img *models.Image) error {
	depSvc := ctx.GetHTTPService("deployment-service")

	body, header, err := getForm(img)
	if err != nil {
		return err
	}

	resp, err := depSvc.PostWithHeaders(ctx, "deploy", nil, body, header)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func getForm(img *models.Image) ([]byte, map[string]string, error) {
	file, err := os.Open(imageZipName)
	if err != nil {
		return nil, nil, err
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	defer writer.Close()

	part, err := writer.CreateFormFile("image", imageZipName)
	if err != nil {
		return nil, nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	err = addField(writer, "name", img.Name)
	err = addField(writer, "tag", img.Tag)
	err = addField(writer, "region", img.Region)
	err = addField(writer, "repository", img.Repository)
	err = addField(writer, "serviceID", img.ServiceID)
	err = addField(writer, "repository", img.Repository)
	err = addField(writer, "region", img.Region)
	err = addField(writer, "loginServer", img.LoginServer)
	err = addField(writer, "serviceName", img.ServiceName)
	err = addField(writer, "accountID", img.AccountID)
	err = addField(writer, "cloudProvider", img.CloudProvider)

	creds, _ := writer.CreateFormField("serviceCreds")
	b, _ := json.Marshal(img.ServiceCreds)
	_, _ = creds.Write(b)

	if err != nil {
		return nil, nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}

	return body.Bytes(), map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}, nil
}

func addField(writer *multipart.Writer, key, value string) error {
	if value == "" {
		return nil
	}

	return writer.WriteField(key, value)
}
