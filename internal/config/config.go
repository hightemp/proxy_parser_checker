package config

import (
	"fmt"
	"os"
	"time"

	y "gopkg.in/yaml.v3"
)

type Config struct {
	SitesForParsing     []string `yaml:"sites_for_parsing"`
	ParsePeriod         string   `yaml:"parse_period"`
	ParsePeriodDuration time.Duration
}

var c Config

func Load(path string) error {
	text, err := os.ReadFile(path)

	if err != nil {
		return fmt.Errorf("Can't read file %s: %v", path, err)
	}

	err = y.Unmarshal(text, &c)

	if err != nil {
		return fmt.Errorf("Can't parse config %s: %v", path, err)
	}

	c.ParsePeriodDuration, err = time.ParseDuration(c.ParsePeriod)

	if err != nil {
		return fmt.Errorf("Can't parse duration in 'ParsePeriod': %v", err)
	}

	return nil
}

func GetConfig() *Config {
	return &c
}
