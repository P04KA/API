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

	// Пытаемся прочитать конфиг, но не падаем если его нет
	if err := viper.ReadInConfig(); err != nil {
		// Просто логируем, но продолжаем
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Конфиг файл не найден - используем значения по умолчанию
		}
	}

	viper.SetDefault("http.port", 8080)
	viper.SetDefault("db.dburl", "postgresql://postgres:postgres@postgres:5432/postgres")

	// Переменные окружения имеют приоритет
	viper.AutomaticEnv()
	viper.SetEnvPrefix("API")
	viper.BindEnv("db.dburl", "DB_URL")

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
