package core

import (
	"fmt"

	"github.com/spf13/viper"
)

// Define all your config in config.yml and add the model here
// Configurations exported
type Configurations struct {
	BOT bot
}

type bot struct {
	TOKEN string
}

// Export the Sugar logger so it can be accessed in other files.
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
