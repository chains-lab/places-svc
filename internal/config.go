package internal

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ServeiceConfig struct {
	Name string `mapstructure:"name"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type RestConfig struct {
	Port     string `mapstructure:"port"`
	Timeouts struct {
		Read       time.Duration `mapstructure:"read"`
		ReadHeader time.Duration `mapstructure:"read_header"`
		Write      time.Duration `mapstructure:"write"`
		Idle       time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeouts"`
}

type DatabaseConfig struct {
	SQL struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"sql"`
}

type JWTConfig struct {
	User struct {
		AccessToken struct {
			SecretKey string `mapstructure:"secret_key"`
		} `mapstructure:"access_token"`
	} `mapstructure:"user"`
	Service struct {
		SecretKey string `mapstructure:"secret_key"`
	} `mapstructure:"service"`
}

type Config struct {
	Service  ServeiceConfig `mapstructure:"service"`
	Log      LogConfig      `mapstructure:"log"`
	Rest     RestConfig     `mapstructure:"rest"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Database DatabaseConfig `mapstructure:"database"`
}

func LoadConfig() (Config, error) {
	configPath := os.Getenv("KV_VIPER_FILE")
	if configPath == "" {
		return Config{}, errors.New("KV_VIPER_FILE env var is not set")
	}
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, errors.Errorf("error reading config file: %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, errors.Errorf("error unmarshalling config: %s", err)
	}

	return config, nil
}
