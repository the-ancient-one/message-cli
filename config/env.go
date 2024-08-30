/*
config provides configuration functions for the message-cli application.
*/
package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// EnvLoad loads the environment variables from the .env file
func EnvLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

// SignMode returns the signature mode from the environment variable file (SignMode)
func SignMode() string {
	EnvLoad()
	return os.Getenv("SignMode")
}

// KemMode returns the KEM mode from the environment variable file (KemMode)
func KemMode() string {
	EnvLoad()
	return os.Getenv("KemMode")
}

// LogFile returns the log file name from the environment variable file (LogFile)
func LogFile() string {
	EnvLoad()
	filename := os.Getenv("LogFile") + "logfile_" + time.Now().Format("2006-01-02")
	return filename
}
