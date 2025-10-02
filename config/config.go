package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP struct {
		Port int
	}
	DB struct {
		DBURL string
	}
}

func LoadConfig(path string) (config Config, err error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, errors.Wrap(err, "cfg not read")
		}
	}

	if err := viper.MergeInConfig(); err != nil {
		return config, errors.Wrap(err, "reading env")
	}

	viper.SetDefault("http.port", 8080)
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		return config, errors.Wrap(err, "unmarshaling cfg")
	}
	return config, nil
}
