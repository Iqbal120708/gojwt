package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    AppPort string
    AppEnv string
    AccessSecret string
    RefreshSecret string
}

func Load() {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system env")
		}

		cfg = &Config{
			DBHost:        os.Getenv("DB_HOST"),
			DBPort:        os.Getenv("DB_PORT"),
			DBUser:        os.Getenv("DB_USER"),
			DBPassword:    os.Getenv("DB_PASSWORD"),
			DBName:        os.Getenv("DB_NAME"),
			AppPort:       os.Getenv("APP_PORT"),
			AppEnv:        os.Getenv("APP_ENV"),
			AccessSecret:  os.Getenv("ACCESS_SECRET_KEY"),
			RefreshSecret: os.Getenv("REFRESH_SECRET_KEY"),
		}
	})
}

func Get() *Config {
	if cfg == nil {
		log.Fatal("Config not loaded. Call config.Load() first")
	}
	return cfg
}