package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func parseConfig(configFile string) (ConfigJSON, error) {
	var config ConfigJSON
	if configFile == "" {
		// default to the normal location
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return config, err
		}
		configFile = filepath.Join(homeDir, ".hitec.conf")
	}

	// Validate if the file exists
	_, err := os.Stat(configFile)
	if err != nil {
		// If the error is due to the file not existing
		if os.IsNotExist(err) {
			// Create the file
			config.InfluxBucket = "Bucket Name"
			config.InfluxOrg = "Org Name"
			config.InfluxToken = "Security token"
			config.ServerAddress = "http://serveraddress:8086"
			outData, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return config, err
			}
			err = ioutil.WriteFile(configFile, outData, 0644)
			if err != nil {
				return config, err
			}
			// Return explanation
			return config, fmt.Errorf("config file not found, generated at %s", configFile)
		}
		return config, err
	}

	// File existed, read and parse it
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}
	// Return the config
	return config, nil
}
