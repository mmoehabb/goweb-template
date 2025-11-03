package constants

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	Port        int
}

var AppConfig = getConfig()

func getConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	return Config{
		DatabaseUrl: dbURL,
		Port:        portNumber,
	}
}
