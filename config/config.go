package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type EnvConfig struct {
	ElasticCloudID string `mapstructure:"ELASTIC_CLOUD_ID"`
	ElasticAPIKey  string `mapstructure:"ELASTIC_API_KEY"`
	TMDBAPIKey     string `mapstructure:"TMDB_API_KEY"`
	TMDBBaseURL    string `mapstructure:"TDMB_BASE_URL"`
}

func (e *EnvConfig) Validate() error {
	if e.ElasticCloudID == "" {
		return errors.New("ELASTIC_CLOUD_ID is not set")
	}
	if e.ElasticAPIKey == "" {
		return errors.New("ELASTIC_API_KEY is not set")
	}
	if e.TMDBAPIKey == "" {
		return errors.New("TMDB_API_KEY is not set")
	}
	if e.TMDBBaseURL == "" {
		return errors.New("TDMB_BASE_URL is not set")
	}
	return nil
}

func LoadEnvConfig(path string) (*EnvConfig, error) {
	fileExists := fileExists(path)
	if !fileExists {
		return nil, fmt.Errorf("%s does not exist", path)
	}

	v := viper.New()
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg EnvConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}
