package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func AesPasswd() string {
	EnvLoad()
	return os.Getenv("aesPasswd")
}
