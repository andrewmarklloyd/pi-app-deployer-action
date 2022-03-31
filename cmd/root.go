package cmd

import (
	"fmt"
	"os"

	"github.com/andrewmarklloyd/pi-app-deployer-action/internal/pkg/config"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "pi-app-deployer-action",
	Short: "Github action that supports deploying applications to Raspberry Pi's.",
	Long:  `Github action that supports deploying applications to Raspberry Pi's. See the pi-app-deployer repo for more information: https://github.com/andrewmarklloyd/pi-app-deployer`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
	// TODO: exit 1 if no sub command ran?
}

func init() {
	rootCmd.PersistentFlags().String("repoName", "", "Name of the Github repo including the owner")
	rootCmd.PersistentFlags().String("manifestName", "", "Name of the pi-app-deployer manifest")
	rootCmd.PersistentFlags().String("host", "", "Name of the pi-app-deployer host")
}

// TODO: is there a better way to get env vars but not expose them in the flag error output?
func envVarsMust() config.EnvVarConfig {
	apiKey := os.Getenv("PI_APP_DEPLOYER_API_KEY")
	if apiKey == "" {
		fmt.Println("PI_APP_DEPLOYER_API_KEY env var is required")
		os.Exit(1)
	}

	githubSHA := os.Getenv("GITHUB_SHA")
	if githubSHA == "" {
		fmt.Println("GITHUB_SHA env var is required")
		os.Exit(1)
	}

	return config.EnvVarConfig{
		APIKey:    apiKey,
		GithubSHA: githubSHA,
	}
}
