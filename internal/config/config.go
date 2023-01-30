// Package config
package config

import (
	"io/ioutil"
)

// C is the global singleton variable for the config.
var C *Config

// Config is the struct for the config file.
type Config struct {
	Server string `yaml:"server"`
	Users  []User `yaml:"users"`
}

// User is the struct for the user.
type User struct {
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

// LoadConfig loads the config file.
func LoadConfig(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &C)
	if err != nil {
		return err
	}

	return nil
}
