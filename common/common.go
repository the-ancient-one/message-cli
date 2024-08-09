package common

import (
	"log/slog"
	"message-cli/config"
	"os"
	"path/filepath"
)

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
