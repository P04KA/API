package config

import (
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

	viper.SetDefault("http.port", 8080)
	viper.SetDefault("db.dburl", "postgresql://postgres:postgres@postgres:5432/postgres")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("API")
	viper.BindEnv("db.dburl", "DB_URL")

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
