package lib

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	// Search config in current directory with name ".log_analyzer" (without extension).
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".log_analyzer")

	viper.AutomaticEnv() // read in environment variables that match

	setupDefaults() // setup defaults config state
	loadConfig()    // load the configuration for the desired file

	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
}

func setupDefaults() {
	viper.SetDefault("scan.interval", 5*time.Second)
	viper.SetDefault("scan.error.path", []string{"C:/Program Files/Apache24/logs/error.log"})
	viper.SetDefault("scan.error.keep_logs", true)
	viper.SetDefault("scan.access.path", []string{"C:/Program Files/Apache24/logs/access.log"})
	viper.SetDefault("scan.access.keep_logs", true)

	viper.SetDefault("database.url", "./logs/log_analyzer.db")
	viper.SetDefault("database.batch_size", 250)

	viper.SetDefault("rotate", true)
}

func loadConfig() {
	// If a config file is found, read it
	if err := viper.ReadInConfig(); err != nil {
		handleConfigError(err)
	}
}

func handleConfigError(err error) {
	// If the config file is not found, then create a default one
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		createDefaultConfig()
		return
	}
	// Error while reading the config file
	log.Fatalf("Error when reading config file -> %v", err)
}

// Create a config file with default options
func createDefaultConfig() {
	if err := viper.SafeWriteConfig(); err != nil {
		log.Fatalf("Error when writing default config -> %v", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error when reading config file -> %v", err)
	}
}
