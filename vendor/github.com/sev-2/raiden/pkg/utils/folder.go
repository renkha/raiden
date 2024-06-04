package utils

import (
	"os"
	"path/filepath"
)

func IsFolderExists(folderPath string) bool {
	_, err := os.Stat(folderPath)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func GetCurrentDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func CreateFolder(folder string) error {
	return os.Mkdir(folder, os.ModePerm)
}

func DeleteFolder(folderPath string) error {
	if !IsFolderExists(folderPath) {
		return nil
	}

	return os.RemoveAll(folderPath)
}

func GetAbsolutePath(path string) (string, error) {
	currDir, err := GetCurrentDirectory()
	if err != nil {
		return "", err
	}

	return filepath.Join(currDir, path), nil
}
