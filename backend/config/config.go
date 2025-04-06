package config

import (
	"encoding/json"
	"os"
	"time"
)

func Init() (err error) {
	fileName := "env/config.json"
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		return
	}
	timeout, err := time.ParseDuration(Config.Timeout)
	if err != nil {
		return
	}
	Config.ReadTimeout = timeout
	Config.WriteTimeout = timeout
	Config.Storage = "storage/upload"
	return
}
