package models

type Deploy struct {
	ClusterName    string      `json:"cluster_name"`
	Region         string      `json:"region"`
	ServiceName    string      `json:"service_name"`
	CloudProvider  string      `json:"cloud_provider"`
	Namespace      string      `json:"namespace"`
	DockerRegistry string      `json:"docker_registry"`
	Key            interface{} `json:"key"`
}

type GoogleCred struct {
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	AuthUri                 string `json:"auth_uri"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
	PrivateKey              string `json:"private_key"`
	PrivateKeyId            string `json:"private_key_id"`
	ProjectId               string `json:"project_id"`
	TokenUri                string `json:"token_uri"`
	Type                    string `json:"type"`
	UniverseDomain          string `json:"universe_domain"`
}
