package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LogFilePath	string `json:"log_file_path"`
}

func GetConfigFromFile() Config {
	c := Config{}

	f, err := os.Open("configapp.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Opening config file: %s\n", err.Error())
		return c
	}

	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Decoding config file: %s\n", err.Error())
		return c
	}

	return c
}