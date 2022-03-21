package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

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

	repoName := flag.String("repoName", "", "Name of the Github repo including the owner")
	manifestName := flag.String("manifestName", "", "Name of the pi-app-deployer manifest")
	flag.Parse()

	if *repoName == "" {
		fmt.Println("repoName is required")
		os.Exit(1)
	}

	if *manifestName == "" {
		fmt.Println("manifestName is required")
		os.Exit(1)
	}

	artifact := Artifact{
		SHA:          githubSha,
		RepoName:     *repoName,
		Name:         fmt.Sprintf("app_%s", githubSha),
		ManifestName: *manifestName,
	}

	err := triggerDeploy(apiKey, artifact)
	if err != nil {
		fmt.Println("Error triggering deploy:", err)
		os.Exit(1)
	}

	err = waitForSuccessfulDeploy(apiKey, artifact)
	if err != nil {
		fmt.Println("Error checking deploy status:", err)
		os.Exit(1)
	}
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
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO: check status fmt.Println(string(body))
	return nil
}

func waitForSuccessfulDeploy(apiKey string, artifact Artifact) error {
	max := 24
	count := 0
	condition := "UNKNOWN"
	for {
		if count >= max {
			return fmt.Errorf("Max number of retries exceeded. Deploy condition: %s", condition)
		}
		if condition == "SUCCESS" {
			break
		}

		fmt.Println(fmt.Sprintf("Attempt number %d", count))
		status, err := checkDeployStatus(apiKey, artifact)
		if err != nil {
			return err
		}

		if status.Status != "success" {
			return fmt.Errorf("Receieved a non successful status '%s' while getting deploy status", status)
		}

		fmt.Println("Deploy condition:", status.Condition)
		count += 1
		time.Sleep(5 * time.Second)
	}
	return nil
}

func checkDeployStatus(apiKey string, artifact Artifact) (DeployStatus, error) {
	j, err := json.Marshal(artifact)
	if err != nil {
		return DeployStatus{}, err
	}

	req, err := http.NewRequest("GET", "https://pi-app-deployer.herokuapp.com/deploy/status", bytes.NewBuffer(j))
	if err != nil {
		return DeployStatus{}, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return DeployStatus{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DeployStatus{}, err
	}
	defer resp.Body.Close()

	deployStatus := DeployStatus{}
	err = json.Unmarshal(body, &deployStatus)
	if err != nil {
		return DeployStatus{}, err
	}

	return deployStatus, nil
}
