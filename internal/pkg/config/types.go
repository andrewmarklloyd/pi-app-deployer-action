package config

type Artifact struct {
	RepoName     string `json:"repoName"`
	ManifestName string `json:"manifestName"`
	SHA          string `json:"sha"`
	Name         string `json:"name"`
}

type DeployStatus struct {
	Status    string `json:"status"`
	Condition string `json:"condition"`
	Error     string `json:"error"`
}

type EnvVarConfig struct {
	APIKey    string
	GithubSHA string
}
