// config/config.go - Add storage configuration
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

	// File Storage Configuration
	StorageProvider   string `mapstructure:"STORAGE_PROVIDER"`     // "local" or "s3"
	StorageBasePath   string `mapstructure:"STORAGE_BASE_PATH"`    // For local storage
	StorageBaseURL    string `mapstructure:"STORAGE_BASE_URL"`     // For local storage
	StorageCDNDomain  string `mapstructure:"STORAGE_CDN_DOMAIN"`   // Optional CDN domain
	
	// S3 Configuration
	S3Region          string `mapstructure:"S3_REGION"`
	S3Bucket          string `mapstructure:"S3_BUCKET"`
	S3AccessKeyID     string `mapstructure:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey string `mapstructure:"S3_SECRET_ACCESS_KEY"`
	
	// File Upload Limits
	MaxFileSize   int64 `mapstructure:"MAX_FILE_SIZE"`     // in bytes
	MaxAvatarSize int64 `mapstructure:"MAX_AVATAR_SIZE"`   // in bytes
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

	// Set defaults
	if config.TokenHourLifespan <= 0 {
		config.TokenHourLifespan = 24
	}
	
	if config.StorageProvider == "" {
		config.StorageProvider = "local"
	}
	
	if config.StorageBasePath == "" {
		config.StorageBasePath = "./uploads"
	}
	
	if config.StorageBaseURL == "" {
		config.StorageBaseURL = "http://localhost:8081/uploads"
	}
	
	if config.MaxFileSize <= 0 {
		config.MaxFileSize = 50 * 1024 * 1024 // 50MB default
	}
	
	if config.MaxAvatarSize <= 0 {
		config.MaxAvatarSize = 10 * 1024 * 1024 // 10MB default
	}

	return
}

func (c *Config) GetTokenDuration() time.Duration {
	return time.Duration(c.TokenHourLifespan) * time.Hour
}