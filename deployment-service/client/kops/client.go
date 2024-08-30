package kops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	gofrSvc "gofr.dev/pkg/gofr/service"

	"kops.dev/deployment-service/models"
)

type client struct {
	kopsSvc gofrSvc.HTTP
}

func New(svc gofrSvc.HTTP) *client {
	return &client{kopsSvc: svc}
}

func (c *client) GetServiceCreds(ctx context.Context, serviceID string) (*models.Credentials, error) {
	var data models.Response

	api := fmt.Sprintf("/service/%s/service-credentials")

	resp, err := c.kopsSvc.Get(ctx, api, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}

	return &data.Data, nil
}

type imageUpdate struct {
	Image string `json:"image"`
}

var errService = errors.New("non OK status code received")

func (c *client) UpdateImage(ctx context.Context, serviceId, imageURL string) error {
	api := fmt.Sprintf("/service/%s/image", serviceId)

	payload, _ := json.Marshal(imageUpdate{Image: imageURL})

	resp, err := c.kopsSvc.Put(ctx, api, nil, payload)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errService
	}

	return nil
}
