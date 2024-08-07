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

func SignMode() string {
	EnvLoad()
	return os.Getenv("SignMode")
}

func KemMode() string {
	EnvLoad()
	return os.Getenv("KemMode")
}
