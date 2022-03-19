package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Running in Go, neat!")
	apiKey := os.Getenv("PI_APP_DEPLOYER_API_KEY")
	if apiKey == "" {
		fmt.Println("PI_APP_DEPLOYER_API_KEY is required")
		os.Exit(1)
	}
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
