package config

import "fmt"

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Addr string `mapstructure:"addr"`
}

// GetAddr get server address
func (c ServerConfig) GetAddr() string {
	if c.Addr != "" {
		return c.Addr
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
