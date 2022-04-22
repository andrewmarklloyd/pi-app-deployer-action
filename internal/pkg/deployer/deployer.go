package deployer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andrewmarklloyd/pi-app-deployer-action/internal/pkg/config"
	"github.com/andrewmarklloyd/pi-app-deployer/api/v1/status"
)

func TriggerDeploy(apiKey, host string, artifact config.Artifact) error {
	j, err := json.Marshal(artifact)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/push", host), bytes.NewBuffer(j))

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

	var apiRes config.APIResponse
	err = json.Unmarshal(body, &apiRes)
	if err != nil {
		return err
	}
	if apiRes.RequestStatus != "success" {
		return fmt.Errorf("deploy status error: %s", apiRes.Error)
	}

	return nil
}

func WaitForSuccessfulDeploy(apiKey, host string, artifact config.Artifact) error {
	max := 24
	count := 0
	status := "UNKNOWN"
	for {
		if count >= max {
			return fmt.Errorf("Max number of retries exceeded. Deploy status: %s", status)
		}
		if status == "SUCCESS" {
			break
		}

		fmt.Println(fmt.Sprintf("Attempt number %d", count))
		cond, err := CheckDeployCondition(apiKey, host, artifact)
		if err != nil {
			return err
		}

		status = cond.Status
		fmt.Println("Deploy status:", status)
		count += 1
		time.Sleep(5 * time.Second)
	}
	return nil
}

func CheckDeployCondition(apiKey, host string, artifact config.Artifact) (status.UpdateCondition, error) {
	cond := status.UpdateCondition{}
	j, err := json.Marshal(artifact)
	if err != nil {
		return cond, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/deploy/status", host), bytes.NewBuffer(j))
	if err != nil {
		return cond, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return cond, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cond, err
	}
	defer resp.Body.Close()

	r := config.APIResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return cond, err
	}

	if r.RequestStatus != "success" {
		return cond, fmt.Errorf("error from api response: %s", r.Error)
	}

	return r.UpdateCondition, nil
}
