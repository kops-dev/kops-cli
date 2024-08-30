package models

import "gofr.dev/pkg/gofr/file"

type Image struct {
	Data       file.Zip `file:"image"`
	Name       string   `form:"name"`
	Tag        string   `form:"tag"`
	ServiceID  string   `form:"serviceID"`
	Repository string   `form:"repository"`
	Region     string   `form:"region"`
}
