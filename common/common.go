package common

import (
	"os"
	"path/filepath"
)

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
