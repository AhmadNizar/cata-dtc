package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Pokemon  PokemonConfig
}

type AppConfig struct {
	Name    string
	Version string
	Host    string
	Port    string
	Env     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type PokemonConfig struct {
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int
	CacheTTL   time.Duration
}

func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "pokemon-api"),
			Version: getEnv("API_VERSION", "v1"),
			Host:    getEnv("APP_HOST", "localhost"),
			Port:    getEnv("APP_PORT", "8080"),
			Env:     getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("MYSQL_HOST", "mysql"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_ROOT_PASSWORD", "password"),
			Name:     getEnv("MYSQL_DATABASE", "pokemon_db"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "redis"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Pokemon: PokemonConfig{
			BaseURL:    getEnv("POKEMON_API_URL", "https://pokeapi.co/api/v2"),
			Timeout:    getEnvAsDuration("POKEMON_API_TIMEOUT", "30s"),
			MaxRetries: getEnvAsInt("POKEMON_API_MAX_RETRIES", 3),
			CacheTTL:   getEnvAsDuration("POKEMON_CACHE_TTL", "5m"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
