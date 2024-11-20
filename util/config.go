package util

import (
	"github.com/spf13/viper"
)

// Configuration holds the application configuration
type Configuration struct {
	AppName string    `mapstructure:"app_name"`
	Port    string    `mapstructure:"port"`
	Debug   bool      `mapstructure:"debug"`
	Jwtkey  string    `mapstructure:"jwtkey"`
	DB      DbConfig  `mapstructure:"db"`
	Dir     DirConfig `mapstructure:"dir"`
}

// DbConfig holds the database configuration
type DbConfig struct {
	Host     string `mapstructure:"host"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type DirConfig struct {
	Uploads string `mapstructure:"uploads"`
	Logs    string `mapstructure:"logs"`
}

// InitConfig initializes and reads configuration using Viper
func InitConfig() (Configuration, error) {
	// Set the file name and type for the .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".") // Look for .env in the current directory

	// Set default values
	viper.SetDefault("app_name", "MyApp")
	viper.SetDefault("port", "8080")
	viper.SetDefault("debug", true)
	viper.SetDefault("jwtkey", "ec0mM3RceAPP")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.name", "ecommerce-db")
	viper.SetDefault("db.username", "postgres")
	viper.SetDefault("db.password", "postgres")
	viper.SetDefault("dir.uploads", "./uploads")
	viper.SetDefault("dir.logs", "./logs")

	// Read the .env file if it exists
	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, err
	}

	// Bind environment variables
	viper.AutomaticEnv()

	// Unmarshal configuration into the Configuration struct
	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}
