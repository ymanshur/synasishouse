package appctx

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ymanshur/synasishouse/order/internal/consts"
)

var (
	oneConfig  *Config
	configOnce sync.Once
)

// Config stores all configurations
type Config struct {
	Environment string `mapstructure:"environment"`

	HTTPServer ServerConfig `mapstructure:"http_server"`

	DB             DBConfig `mapstructure:"db"`
	DBMigrationURL string   `mapstructure:"db_migration_url"`

	GRPCClient GRPCClient `mapstructure:"grpc_client"`
}

type DBConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"pass"`
}

// GRPCClient stores all gRPC client configurations
type GRPCClient struct {
	Inventory ClientConfig `mapstructure:"inventory"`
}

type ClientConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Addr            string `mapstructure:"addr"`
	MaxRetry        uint
	PerRetryTimeout time.Duration
}

// GetAddr get client address
func (c ClientConfig) GetAddr() string {
	if c.Addr != "" {
		return c.Addr
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

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

// LoadConfig return config instance.
// It will read [consts.DefaultConfigFilename] file with [consts.DefaultConfigExt] extension
func LoadConfig() *Config {
	configOnce.Do(func() {
		var err error
		oneConfig, err = LoadConfigWithFilename(consts.DefaultConfigFilename, consts.DefaultConfigExt)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot load config")
		}
	})

	return oneConfig
}

// LoadConfigWithFilename reads configuration from a given filename
// at root project directory or environment variables.
func LoadConfigWithFilename(filename, ext string) (*Config, error) {
	path := filepath.Join(rootDir())
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType(ext) // json, yml, etc.

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	// Replace dots in keys with underscores for environment variable compatibility
	// (e.g., "db.host" becomes "DB_HOST")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	var config Config

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &config, nil
}

// rootDir get an absolute root dir of current project
func rootDir() string {
	_, b, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Join(filepath.Dir(b), "..", "..")
	}
	return ""
}
