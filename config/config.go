package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	General General `toml:"general"`
}

type General struct {
	Port    string `toml:"port"`
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Debug   bool   `toml:"debug"`
}

// LoadConfig deserialize config from toml file.
func LoadConfig(fp string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(fp, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
