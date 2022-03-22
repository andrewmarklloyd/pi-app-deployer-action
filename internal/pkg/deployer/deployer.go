package deployer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andrewmarklloyd/pi-app-deployer-action/internal/pkg/config"
)

func TriggerDeploy(apiKey string, artifact config.Artifact) error {
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

	var deployStatus config.DeployStatus
	err = json.Unmarshal(body, &deployStatus)
	if err != nil {
		return err
	}
	if deployStatus.Error != "" {
		return fmt.Errorf("deploy status error: %s", deployStatus.Error)
	}

	return nil
}

func WaitForSuccessfulDeploy(apiKey string, artifact config.Artifact) error {
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
		status, err := CheckDeployStatus(apiKey, artifact)
		if err != nil {
			return err
		}

		if status.Status != "success" {
			return fmt.Errorf("Receieved a non successful status '%s' while getting deploy status", status)
		}

		condition = status.Condition
		fmt.Println("Deploy condition:", condition)
		count += 1
		time.Sleep(5 * time.Second)
	}
	return nil
}

func CheckDeployStatus(apiKey string, artifact config.Artifact) (config.DeployStatus, error) {
	deployStatus := config.DeployStatus{}
	j, err := json.Marshal(artifact)
	if err != nil {
		return deployStatus, err
	}

	req, err := http.NewRequest("GET", "https://pi-app-deployer.herokuapp.com/deploy/status", bytes.NewBuffer(j))
	if err != nil {
		return deployStatus, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return deployStatus, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return deployStatus, err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &deployStatus)
	if err != nil {
		return deployStatus, err
	}

	return deployStatus, nil
}
