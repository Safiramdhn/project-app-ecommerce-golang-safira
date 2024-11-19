package util

import "github.com/spf13/viper"

type Configuration struct {
	AppName string
	Port    string
	Debug   bool
	DB      DbConfig
}

type DbConfig struct {
	Host     string
	Name     string
	Username string
	Password string
}

func ReadConfigurations() (Configuration, error) {
	viper.AutomaticEnv()

	return Configuration{
		AppName: viper.GetString("APP_NAME"),
		Port:    viper.GetString("PORT"),
		Debug:   viper.GetBool("DEBUG"),
		DB: DbConfig{
			Host:     viper.GetString("DB_HOST"),
			Name:     viper.GetString("DB_NAME"),
			Username: viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
		},
	}, nil
}
