package config

import (
	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	Address        string `toml:"address"`
	Port           int    `toml:"port"`
	MaxConnections int    `toml:"max_connections"`
	Weight         int    `toml:"weight"`
}

type LoadBalancer struct {
	Address   string         `toml:"address"`
	Port      int            `toml:"port"`
	Algorithm string         `toml:"algorithm"`
	Layer     int            `toml:"layer"`
	Servers   []ServerConfig `toml:"server"`
}

type Config struct {
	LoadBalancer LoadBalancer `toml:"load_balancer"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
