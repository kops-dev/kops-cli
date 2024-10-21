package deploy

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"gofr.dev/pkg/gofr"

	zopClient "zop.dev/client"
	"zop.dev/models"
)

var (
	errUpdatingImage = errors.New("unable to update the image for your service via zop.dev services")
)

type client struct {
}

func New() zopClient.ServiceDeployer {
	return &client{}
}

func (*client) Deploy(ctx *gofr.Context, img *models.Image, zipFile string) error {
	depSvc := ctx.GetHTTPService("deployment-service")

	body, header, err := getForm(img, zipFile)
	if err != nil {
		return err
	}

	resp, err := depSvc.PostWithHeaders(ctx, "deploy", nil, body, header)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		ctx.Logger.Errorf("error communicating with the deployment service, status code returned - %d", resp.StatusCode)

		return errUpdatingImage
	}

	return nil
}

func getForm(img *models.Image, zipFile string) (bodyBytes []byte, headers map[string]string, err error) {
	file, err := os.Open(zipFile)
	if err != nil {
		return nil, nil, err
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	defer writer.Close()

	part, err := writer.CreateFormFile("image", zipFile)
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
