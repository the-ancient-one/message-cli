package config

import (
	"log"
	"os"
	"time"

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

func LogFile() string {
	EnvLoad()
	filename := os.Getenv("LogFile") + "logfile_" + time.Now().Format("2006-01-02")
	return filename
}
