package cmd

import (
	"fmt"
	"os"

	"github.com/andrewmarklloyd/pi-app-deployer-action/internal/pkg/config"
	"github.com/andrewmarklloyd/pi-app-deployer-action/internal/pkg/deployer"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "TODO",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		runDeploy(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) {
	envVarConfig := envVarsMust()
	repoName, err := cmd.Flags().GetString("repoName")
	if err != nil {
		fmt.Println("error getting repoName flag", err)
		os.Exit(1)
	}
	if repoName == "" {
		fmt.Println("repoName flag is required")
		os.Exit(1)
	}

	manifestName, err := cmd.Flags().GetString("manifestName")
	if err != nil {
		fmt.Println("error getting repoName flag", err)
		os.Exit(1)
	}
	if manifestName == "" {
		fmt.Println("manifestName flag is required")
		os.Exit(1)
	}

	artifact := config.Artifact{
		SHA:          envVarConfig.GithubSHA,
		RepoName:     repoName,
		Name:         fmt.Sprintf("app_%s", envVarConfig.GithubSHA),
		ManifestName: manifestName,
	}

	err = deployer.TriggerDeploy(envVarConfig.APIKey, artifact)
	if err != nil {
		fmt.Println("Error triggering deploy:", err)
		os.Exit(1)
	}
	fmt.Println("Successfully triggered deploy, waiting for successful deployment.")

	err = deployer.WaitForSuccessfulDeploy(envVarConfig.APIKey, artifact)
	if err != nil {
		fmt.Println("Error checking deploy status:", err)
		os.Exit(1)
	}
}
