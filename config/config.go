package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LogFilePath	string 		`json:"log_file_path"`
	DB			DBConfig 	`json:"db"`
}

type DBConfig struct {
	Host		string 	`json:"host"`
	Port		int		`json:"port"`
	User		string	`json:"user"`
	Password	string	`json:"password"`
	Name		string	`json:"name"`
	SSL			bool	`json:"ssl"`
}

func GetConfigFromFile() Config {
	c := Config{}

	f, err := os.Open("configapp.json")
	defer f.Close()

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