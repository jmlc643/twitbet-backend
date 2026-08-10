package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	JWTSecret     string
	FrontendURL   string
	CloudinaryURL string
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPass      string
	SMTPSender    string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: No se encontró archivo .env local. Se usarán las variables del entorno del sistema.")
	}

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisURL:    getEnv("REDIS_URL", ""),
		JWTSecret:     getEnv("JWT_SECRET", "default_secret"),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:5173"),
		CloudinaryURL: getEnv("CLOUDINARY_URL", ""),
		SMTPHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPass:      getEnv("SMTP_PASS", ""),
		SMTPSender:    getEnv("SMTP_SENDER", "noreply@twitbet.com"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("La variable de entorno DATABASE_URL está vacía.")
	}

	if cfg.RedisURL == "" {
		log.Fatal("La variable de entorno REDIS_URL está vacía.")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
