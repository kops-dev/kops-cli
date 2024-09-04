package models

import "gofr.dev/pkg/gofr/file"

type Image struct {
	Data          file.Zip `file:"image"`
	Name          string   `form:"name" json:"name"`
	Tag           string   `form:"tag" json:"tag"`
	ServiceID     string   `form:"serviceID" json:"serviceID"`
	ServiceCreds  any      `form:"serviceCreds" json:"serviceCreds"`
	Repository    string   `form:"repository" json:"repository"`
	Region        string   `form:"region" json:"region"`
	LoginServer   string   `form:"loginServer" json:"loginServer"`
	ServiceName   string   `form:"serviceName" json:"serviceName"`
	AccountID     string   `form:"accountID" json:"accountID"`
	CloudProvider string   `form:"cloudProvider" json:"cloudProvider"`
}
