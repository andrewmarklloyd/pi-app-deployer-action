package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Artifact struct {
	RepoName     string `json:"repoName"`
	ManifestName string `json:"manifestName"`
	SHA          string `json:"sha"`
	Name         string `json:"name"`
}

func main() {
	apiKey := os.Getenv("PI_APP_DEPLOYER_API_KEY")
	if apiKey == "" {
		fmt.Println("PI_APP_DEPLOYER_API_KEY is required")
		os.Exit(1)
	}

	githubSha := os.Getenv("GITHUB_SHA")
	if githubSha == "" {
		fmt.Println("GITHUB_SHA is required")
		os.Exit(1)
	}

	fmt.Println(os.Args)

	artifact := Artifact{
		SHA:          githubSha,
		RepoName:     "",
		Name:         fmt.Sprintf("app_%s", githubSha),
		ManifestName: "",
	}
	fmt.Println(artifact)
	// err := triggerDeploy(apiKey, payload)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}

func triggerDeploy(apiKey string, artifact Artifact) error {
	j, err := json.Marshal(artifact)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://pi-app-deployer.herokuapp.com/push", bytes.NewBuffer(j))

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(string(body))
	return nil
}

/*

get inputs:
	PI_APP_DEPLOYER_API_KEY
	GITHUB_SHA
	optional timeout
pass payload to api
	https://pi-app-deployer.herokuapp.com/push
if status is not success fail
while not timeout
	curl deploy/status, check for success
*/
