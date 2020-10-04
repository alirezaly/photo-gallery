package config

import "github.com/BurntSushi/toml"

type Config struct {
	Port string `toml:"port"`
	Host string `toml:"host"`
	Database string `toml:"database"`
}

func SetupConfig() (*Config, error) {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
