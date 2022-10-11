package config

import (
	"path/filepath"
	"strings"

	"github.com/shuvava/treehub/internal/utils/fshelper"

	"github.com/fsnotify/fsnotify"
	"github.com/shuvava/go-logging/logger"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = "appsettings"
	defaultConfigPath = "."
	prodConfigName    = "config_prod"
	prodConfigPath    = "./config"
)

// StorageConfig file storage configuration
type StorageConfig struct {
	Type string `mapstructure:"type"`
	Root string `mapstructure:"root"`
}

// DbConfig service database configuration
type DbConfig struct {
	Type             string `mapstructure:"type"`
	ConnectionString string `mapstructure:"connectionString"`
}

// AppConfig root app config
type AppConfig struct {
	Port     int      `mapstructure:"port"`
	LogLevel string   `mapstructure:"logLevel"`
	Db       DbConfig `mapstructure:"db"`

	Storage StorageConfig `mapstructure:"storage"`
}

// OnConfigChange callback for config changes
type OnConfigChange func(*AppConfig)

// NewConfig loads new config or run panic on error
func NewConfig(logger logger.Logger, fn OnConfigChange) *AppConfig {
	log := logger.SetOperation("config-initialization")
	var cfg AppConfig
	if fn != nil {
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.WithField("file", e.Name).
				Debug("Config file was changed")
			if err := viper.Unmarshal(&cfg); err != nil {
				log.WithError(err).
					Error("Error on config Unmarshal:")
			} else {
				log.Info("config auto reload!")
				fn(&cfg)
			}
		})
	}

	// add config with default values
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(defaultConfigPath)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.WithError(err).
			Fatal("Fatal error on reading config file")
	}
	// watching if default config changed
	viper.WatchConfig()

	if fshelper.IsPathExist(filepath.Join(prodConfigPath, prodConfigName)) {
		// add production config
		viper.SetConfigName(prodConfigName)
		viper.AddConfigPath(prodConfigPath)
		// merge values into default one
		err = viper.MergeInConfig()
		if err != nil { // Handle errors reading the config file
			log.WithError(err).
				Fatal("Fatal error on reading config file")
		}
		// watch production config changes
		viper.WatchConfig()
	}

	// take into account env variables with the highest priority
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.WithError(err).
			Fatal("Fatal error on config Unmarshal")
	}

	return &cfg
}

// PrintConfig returns print current config into log output
func (cfg *AppConfig) PrintConfig(log logger.Logger) {
	log.Info("Current config:")
	log.Info("    Port         :", cfg.Port)
	log.Info("    LogLevel     :", cfg.LogLevel)
	log.Info("    Db.Type      :", cfg.Db.Type)
	log.Info("    Storage.Type :", cfg.Storage.Type)
	log.Info("    Storage.Root :", cfg.Storage.Root)
}
