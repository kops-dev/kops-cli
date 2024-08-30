package models

type GCPCreds struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

type AzureCreds struct {
	AppID          string `json:"appId"`
	Password       string `json:"password"`
	ClientSecret   string `json:"clientSecret"`
	ClientID       string `json:"clientId"`
	TenantID       string `json:"tenantId"`
	SubscriptionID string `json:"subscriptionId"`
}

type AWSCreds struct {
	AccessKey    string `json:"aws_access_key_id"`
	AccessSecret string `json:"aws_secret_access_key"`
}

type Credentials struct {
	CloudPlatform  string      `json:"cloudPlatform"`
	ServiceAccID   string      `json:"serviceAccId"`
	ServiceAccCred interface{} `json:"serviceAccCred"`
	Credentials    bool        `json:"credentials"`
}

type Response struct {
	Data Credentials `json:"data"`
}
