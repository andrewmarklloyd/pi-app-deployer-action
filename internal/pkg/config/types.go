package config

import "github.com/andrewmarklloyd/pi-app-deployer/api/v1/status"

type Artifact struct {
	RepoName     string `json:"repoName"`
	ManifestName string `json:"manifestName"`
	SHA          string `json:"sha"`
	Name         string `json:"name"`
}

type EnvVarConfig struct {
	APIKey    string
	GithubSHA string
}

type APIResponse struct {
	RequestStatus     string                            `json:"request"`
	Error             string                            `json:"error"`
	SuccessfulHosts   map[string]status.UpdateCondition `json:"successfulHosts"`
	UnsuccessfulHosts map[string]status.UpdateCondition `json:"unsuccessfulHosts"`
}
