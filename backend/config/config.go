package config

import (
	"encoding/json"
	"os"
	"time"
)
// Init initializes the application configuration by reading from a JSON file,
// parsing the configuration settings, and setting up timeout values.
// It returns an error if any step fails, otherwise it returns nil.
func Init() (err error) {
	// Define the path to the configuration file.
	fileName := "env/config.json"
	// Open the configuration file.
	file, err := os.Open(fileName)
	if err != nil {
		// Return an error if the file cannot be opened.
		return
	}
	// Ensure the file is closed when the function exits.
	defer file.Close()

	// Decode the content of the JSON file into the Config object.
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		// Return an error if the JSON decoding fails.
		return
	}

	// Parse the Timeout field from the configuration and convert it to a time.Duration.
	timeout, err := time.ParseDuration(Config.Timeout)
	if err != nil {
		// Return an error if the timeout parsing fails.
		return
	}

	// Set the read and write timeouts in the Config object.
	Config.ReadTimeout = timeout
	Config.WriteTimeout = timeout

	// Set the storage directory for uploaded files.
	Config.Storage = "storage/upload"
	return
}
