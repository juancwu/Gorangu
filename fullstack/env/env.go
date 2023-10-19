package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB_URL string
var DB_AUTH_TOKEN string

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}
	fmt.Println("Loaded ENV")
	DB_URL = os.Getenv("DB_URL")
	DB_AUTH_TOKEN = os.Getenv("DB_AUTH_TOKEN")
	fmt.Printf("DB_URL: %s\n", DB_URL)
	fmt.Printf("DB_AUTH_TOKEN: %s\n", DB_AUTH_TOKEN)
}
