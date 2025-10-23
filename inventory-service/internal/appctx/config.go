package appctx

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ymanshur/synasishouse/inventory/internal/consts"
)

var (
	oneConfig  *Config
	configOnce sync.Once
)

// Config stores all configurations
type Config struct {
	GRPCServer     ServerConfig `mapstructure:"grpc_server"`
	DB             DBConfig     `mapstructure:"db"`
	DBMigrationURL string       `mapstructure:"db_migration_url"`
}

type DBConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"pass"`
}

// GetURL get database DSN
func (c DBConfig) GetURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
}

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
