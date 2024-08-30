package deploy

import (
	"bytes"
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

	err = getFormField(writer, "name", img.Name)
	err = getFormField(writer, "tag", img.Tag)
	err = getFormField(writer, "region", img.Region)
	err = getFormField(writer, "repository", img.Repository)
	err = getFormField(writer, "serviceID", img.ServiceID)

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

func getFormField(writer *multipart.Writer, key, value string) error {
	k, err := writer.CreateFormField(key)
	if err != nil {
		return err
	}

	_, err = k.Write([]byte(value))
	if err != nil {
		return err
	}

	return nil
}
