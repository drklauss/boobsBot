package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var cfg Config

// Config holds main app config
type Config struct {
	LogPath   string `yaml:"log"`
	*Telegram `yaml:"telegram"`
	*Reddit   `yaml:"reddit"`
	*Imgur    `yaml:"imgur"`
}

// Load reads .yml file and loads it
func Load() error {
	b, err := ioutil.ReadFile("../boobsBot/config.yml")
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return err
	}

	return nil
}

// Get returns config instance
func Get() *Config {
	return &cfg
}
