package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBConnection = ""
	WebPort      = "8200"
	SecretKey    []byte
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBConnection = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=true&loc=Local",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	WebPort = os.Getenv("WEB_PORT")
	SecretKey = []byte(os.Getenv("JWT_SECRET"))
}
