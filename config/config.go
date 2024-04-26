package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBPassword string `mapstructure:"DB_PASSWORD"`
	Limit      int    `mapstructure:"LIMIT"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg *Config
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	
	viper.SetDefault("LIMIT", 10)
	
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	
	return cfg, err
}
