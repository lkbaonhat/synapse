package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration loaded from env / .env file.
type Config struct {
	Port string `mapstructure:"PORT"`
	Env  string `mapstructure:"ENV"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	JWTSecret           string `mapstructure:"JWT_SECRET"`
	JWTAccessTTLMinutes int    `mapstructure:"JWT_ACCESS_TTL_MINUTES"`
	JWTRefreshTTLDays   int    `mapstructure:"JWT_REFRESH_TTL_DAYS"`

	UploadDir string `mapstructure:"UPLOAD_DIR"`

	CORSAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

// AllowedOrigins returns CORS origins as a slice.
func (c *Config) AllowedOrigins() []string {
	if c.CORSAllowedOrigins == "" {
		return []string{"http://localhost:5173"}
	}
	return strings.Split(c.CORSAllowedOrigins, ",")
}

// DSN returns the PostgreSQL data source name.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// Load reads configuration from .env then environment variables.
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("JWT_ACCESS_TTL_MINUTES", 15)
	viper.SetDefault("JWT_REFRESH_TTL_DAYS", 7)
	viper.SetDefault("UPLOAD_DIR", "./uploads")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:5173")

	// Optional .env file — ignore if not present
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal failed: %w", err)
	}
	return &cfg, nil
}
