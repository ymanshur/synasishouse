package config

import "fmt"

type ServerConfig struct {
	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`
	Addr string `mapstructure:"ADDR"`
}

// GetAddr get server address
func (c ServerConfig) GetAddr() string {
	if c.Addr != "" {
		return c.Addr
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
