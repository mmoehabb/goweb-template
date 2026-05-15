package constants

import (
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
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	portNumber, _ := strconv.Atoi(port)
	if portNumber == 0 {
		portNumber = 3000
	}

	return Config{
		DatabaseUrl: dbURL,
		Port:        portNumber,
	}
}
