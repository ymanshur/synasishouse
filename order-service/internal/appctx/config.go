package appctx

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ymanshur/synasishouse/order/internal/consts"
)

var (
	config     Config
	configOnce sync.Once
)

// Config stores all configurations
type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`

	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`

	DBSource       string `mapstructure:"DB_SOURCE"`
	DBMigrationURL string `mapstructure:"DB_MIGRATION_URL"`

	GRPCClientHostInventory string `mapstructure:"GRPC_CLIENT_HOST_INVENTORY"`
	GRPCClientPortInventory int    `mapstructure:"GRPC_CLIENT_PORT_INVENTORY"`
}

type GRPCAddr struct {
	Host string `yaml:"host" json:"host"`
	Port int    `yaml:"port" json:"port"`
}

// LoadConfig load single config instance.
// It will read app.yaml in config directory.
func LoadConfig() Config {
	return LoadConfigWithFilename(consts.DefaultConfigFilename)
}

// LoadConfigWithFilename reads configuration from a given filename or environment variables once
func LoadConfigWithFilename(finename string) Config {
	configOnce.Do(func() {
		path := filepath.Join(rootDir())
		viper.AddConfigPath(path)
		viper.SetConfigName(finename)
		viper.SetConfigType("env") // json, yml, etc.

		// AutomaticEnv will override config file
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal().Err(err)
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatal().Err(err)
		}
	})

	return config
}

// rootDir get an absolute root dir of current project
func rootDir() string {
	_, b, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Join(filepath.Dir(b), "..", "..")
	}
	return ""
}
