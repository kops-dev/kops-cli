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

type GCPInfo struct {
	ProjectID string `json:"project_id"`
}
