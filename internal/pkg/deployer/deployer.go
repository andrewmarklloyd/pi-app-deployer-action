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
	for {
		fmt.Println(fmt.Sprintf("Attempt number %d", count))
		c, err := CheckDeployCondition(apiKey, host, artifact)
		if err != nil {
			return err
		}

		if len(c.UnsuccessfulHosts) == 0 {
			return nil
		}

		count += 1
		time.Sleep(5 * time.Second)
		if count >= max {
			j, _ := json.Marshal(c)
			return fmt.Errorf("Max number of retries exceeded. Deploy conditions from server: %s", j)
		}
	}
}

func CheckDeployCondition(apiKey, host string, artifact config.Artifact) (config.APIResponse, error) {
	c := config.APIResponse{}
	j, err := json.Marshal(artifact)
	if err != nil {
		return c, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/deploy/status", host), bytes.NewBuffer(j))
	if err != nil {
		return c, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()

	r := config.APIResponse{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return c, err
	}

	if r.RequestStatus != "success" {
		return c, fmt.Errorf("error from api response: %s", r.Error)
	}

	return r, nil
}
