package configs

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"server"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	viper.SetConfigFile("configs/config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configs file, %s", err)
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
		return nil, err
	}
	return &cfg, nil
}
