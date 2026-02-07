package main

import (
	"log"
	"os"

	"samokat-scraper/internal/scraper"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env не найден")
	}

	cfg := scraper.Config{
		APIURL:      os.Getenv("API_URL"),
		CategoryURL: os.Getenv("CATEGORY_URL"),
		OutDir:      "data",
		Proxy:       os.Getenv("PROXY"),
		AuthHeader:  getEnv("AUTH_TOKEN", ""),
	}

	if cfg.APIURL == "" {
		log.Fatal("Укажи API_URL в .env")
	}

	if cfg.AuthHeader == "" {
		log.Fatal("Укажи AUTH_TOKEN в .env")
	}

	if err := scraper.Run(cfg); err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	log.Println("Готово!")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
