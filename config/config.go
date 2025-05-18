package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var configPath = []string{
	"./",
	"./config/",
	"./.config/",
	"../.config/",
	"/var/config/application/",
}

type Config struct {
	App struct {
		Env      string
		Port     string
		Log      string
		Database string
		Schema   string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
	}

	RSA struct {
		Private string
		Public  string
	}
}

var config *Config

// NewConfig is a constructor for Config
func new() *Config {
	var err error
	for _, c := range configPath {
		err = godotenv.Load(c + ".env")
		if err != nil {
			log.Printf("failed to laod config from path %s.env", c)
			continue
		}
		break
	}
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config = &Config{}

	config.Redis.Host = os.Getenv("REDIS_HOST")
	config.Redis.Port = os.Getenv("REDIS_PORT")
	config.Redis.Password = os.Getenv("REDIS_PASSWORD")

	config.RSA.Private = os.Getenv("RSA_PRIVATE_KEY")
	config.RSA.Public = os.Getenv("RSA_PUBLIC_KEY")

	config.App.Env = os.Getenv("APP_ENV")
	config.App.Port = os.Getenv("APP_PORT")
	config.App.Log = os.Getenv("APP_LOG")
	config.App.Database = os.Getenv("APP_DATABASE")
	config.App.Schema = os.Getenv("APP_SCHEMA")
	return config
}

func GetConfig() *Config {
	if config == nil {
		config = new()
	}

	return config
}
