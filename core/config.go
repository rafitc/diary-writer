package core

import (
	"fmt"

	"github.com/spf13/viper"
)

// Define all your config in config.yml and add the model here
// Configurations exported
type Configurations struct {
	BOT      bot
	DATABASE DB
	PUBLISH  Publish
	GEN_AI   GEN_AI
}

type bot struct {
	TOKEN string
}

type GEN_AI struct {
	GROK_API_KEY string
	MODEL_NAME   string
}

type DB struct {
	NAME string
	URL  string
}

type Publish struct {
	PUBLISH_JOB_CRON    string
	GITHUB_USERNAME     string
	GITHUB_AUTH_TOKEN   string
	GITHUB_REPO         string
	GITHUB_BRANCH       string
	CLONE_DIRECTORY     string
	GITHUB_COMMIT_USER  string
	GITHUB_COMMIT_EMAIL string
}

var Config *Configurations

func init() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath("config/")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	var configuration *Configurations

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	Config = configuration // Save the config into a global var

}
