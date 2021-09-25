package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// PathExists returns true whether the path exists or not in the file system
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateDir creates a new directory in the file system.
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

func Move(from, to string) error {
	return os.Rename(from, to)
}

func NewFile(fileName string) (*os.File, error) {
	return os.Create(fileName)
}

// WriteToFile writes bytes to a file.
func WriteToFile(fileName string, bytes []byte) error {
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		return err
	}
	_, errorWrite := f.Write(bytes)
	if errorWrite != nil {
		return errorWrite
	}

	return nil
}

// HomeDir returns the current user's home directory
func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		panic("Could not obtain home directory, so the configuration can be created. Error: " + err.Error())
	}

	return usr.HomeDir
}

// DeleteAtPath tries to delete a file given a path
func DeleteAtPath(path string) error {
	return os.Remove(path)
}

// DeleteDir deletes a directory
func DeleteDir(path string) error {
	return os.RemoveAll(path)
}

// IsExtension checks if a file has the correct extension. It assumes the path exists
func IsExtension(path, wantedExtension string) bool {
	ext := filepath.Ext(path)
	return ext == fmt.Sprintf(".%s", wantedExtension)
}
