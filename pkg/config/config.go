package config

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	GracefulTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	// Setup command line flags
	pflag.String("server.host", "0.0.0.0", "Server host")
	pflag.String("server.port", "8080", "Server port")
	pflag.Duration("server.read_timeout", 15*time.Second, "Server read timeout")
	pflag.Duration("server.write_timeout", 15*time.Second, "Server write timeout")
	pflag.Duration("server.graceful_timeout", 15*time.Second, "Server graceful shutdown timeout")
	pflag.Parse()

	// Setup Viper for environment variables
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Bind environment variables
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSLMODE")

	// Set default values
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "bank_db")
	viper.SetDefault("DB_SSLMODE", "disable")

	// Bind command line flags to viper
	viper.BindPFlags(pflag.CommandLine)

	// Server config
	config.Server = ServerConfig{
		Host:            viper.GetString("server.host"),
		Port:            viper.GetString("server.port"),
		ReadTimeout:     viper.GetDuration("server.read_timeout"),
		WriteTimeout:    viper.GetDuration("server.write_timeout"),
		GracefulTimeout: viper.GetDuration("server.graceful_timeout"),
	}

	// Database config
	config.Database = DatabaseConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("DB_NAME"),
		SSLMode:  viper.GetString("DB_SSLMODE"),
	}

	return config, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
