package config

import (
	"log"

	"github.com/spf13/viper"
)

var LocalConfig *Config

type Config struct {
	RABBITMQ_URL string `mapstructure:"RABBITMQ_URL"`
	CPU          int    `mapstructure:"CPU"`
	MemorySize   int    `mapstructure:"MEMORY_SIZE"`
}

func InitConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	var config *Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error reading env file", err)
	}

	return config
}
func SetConfig() {
	LocalConfig = InitConfig()
}
