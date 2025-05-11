package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)


type Config struct {
    Environment string `mapstructure:"ENVIRONMENT"`
    Host        string `mapstructure:"HOST"`
    Port        string `mapstructure:"PORT"`

    DBUsername string `mapstructure:"DB_USERNAME"`
    DBPassword string `mapstructure:"DB_PASSWORD"`
    DBHost     string `mapstructure:"DB_HOSTNAME"`
    DBPort     string `mapstructure:"DB_PORT"`
    DBName     string `mapstructure:"DB_NAME"`

    JWTSecret         string `mapstructure:"JWT_SECRET"`
    TokenHourLifespan int    `mapstructure:"TOKEN_HOUR_LIFESPAN"`
    APISecret         string `mapstructure:"API_SECRET"`

    Version string `mapstructure:"VERSION"`

    RedisHost     string `mapstructure:"REDIS_HOST"`
    RedisPort     string `mapstructure:"REDIS_PORT"`
    RedisPassword string `mapstructure:"REDIS_PASSWORD"`
    RedisDB       int    `mapstructure:"REDIS_DB"`
}

func LoadConfig(name string, path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("config: %v", err)
		return
	}

	if config.TokenHourLifespan <= 0 {
		config.TokenHourLifespan = 24
	}

	return
}

func (c *Config) GetTokenDuration() time.Duration {
	return time.Duration(c.TokenHourLifespan) * time.Hour
}
