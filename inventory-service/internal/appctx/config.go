package appctx

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/ymanshur/synasishouse/inventory/internal/consts"
)

var (
	config     Config
	configOnce sync.Once
)

// Config stores all configurations
type Config struct {
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`

	DBSource       string `mapstructure:"DB_SOURCE"`
	DBMigrationURL string `mapstructure:"DB_MIGRATION_URL"`
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
