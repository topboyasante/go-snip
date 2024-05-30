package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ServerPort string

	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTKey       string
	SMTPUsername string
	SMTPPassword string
	SMTPHost     string
	SMTPAddress  string
}

var ENV = initConfig()

func initConfig() Config {

	godotenv.Load()

	return Config{
		ServerPort:   getEnv("SERVER_PORT", ":4000"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "mypassword"),
		DBName:       getEnv("DB_NAME", "invxice"),
		JWTKey:       getEnv("JWT_KEY", "someJWTKey"),
		SMTPUsername: getEnv("SMTP_USERNAME", "someEmail"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "somePassword"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.emailprovider.com"),
		SMTPAddress:  getEnv("SMTP_ADDR", "smtp.emailprovider.com:someNumber"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
