package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ConfigEnv struct {
	Port               string
	MongoURI           string
	MongoDatabase      string
	AccessTokenSecret  string
	RefreshTokenSecret string
}

func LoadConfig() *ConfigEnv {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	return &ConfigEnv{
		Port:               getEnv("PORT", "3000"),
		MongoURI:           getEnv("MONGO_URI", ""),
		MongoDatabase:      getEnv("MONGO_DATABASE", ""),
		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", ""),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
