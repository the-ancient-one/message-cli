/*
common provides common functions for the message-cli application.
*/
package common

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/the-ancient-one/message-cli/config"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

// SetupLogger sets up the logger for the application
func SetupLogger() *slog.Logger {
	var logfile = config.LogFile()
	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	var handlerOpts = &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}

	logger := slog.New(slog.NewJSONHandler(file, handlerOpts))

	return logger
}

// ListEncryptedMsgFiles lists all the encrypted message files in the user's directory
func ListEncryptedMsgFiles(userID string) ([]string, error) {
	messagesFolder := "storage/" + userID + "/messages/"
	pattern := "encryptedMsg*"

	var encryptedMsgFiles []string

	err := filepath.Walk(messagesFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		matched, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return err
		}
		if !info.IsDir() && matched {
			encryptedMsgFiles = append(encryptedMsgFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return encryptedMsgFiles, nil
}

// CheckUserExists checks if the user exists in the storage directory
func CheckUserExists(directoryPath string) bool {
	directoryPath = filepath.Clean("storage/" + directoryPath)
	fileInfo, err := os.Stat(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	return fileInfo.IsDir()
}

// GetSystemStats returns the system stats for the application
func GetSystemStats() (*memory.Stats, *cpu.Stats, error) {
	mem, err := memory.Get()
	if err != nil {
		return nil, nil, err
	}

	cpu, err := cpu.Get()
	if err != nil {
		return nil, nil, err
	}

	return mem, cpu, nil
}
