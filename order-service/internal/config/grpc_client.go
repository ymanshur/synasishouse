package config

import (
	"fmt"
	"time"
)

// GRPCClient stores all gRPC client configurations
type GRPCClient struct {
	Inventory GRPCClientConfig `mapstructure:"inventory"`
}

type GRPCClientConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Addr            string `mapstructure:"addr"`
	MaxRetry        uint   `mapstructure:"max_retry"`
	PerRetryTimeout time.Duration
}

// GetAddr get client address
func (c GRPCClientConfig) GetAddr() string {
	if c.Addr != "" {
		return c.Addr
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
