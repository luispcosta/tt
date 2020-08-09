package utils

import (
	"fmt"
	"os"
	"os/user"
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

// WriteToFile writes bytes to a file.
func WriteToFile(fileName string, bytes []byte) error {
	fmt.Println(fileName)
	f, err := os.Create(fileName)
	fmt.Println(err)
	defer f.Close()
	if err != nil {
		fmt.Println("Could not create path")
		fmt.Println(err)
		return err
	}
	fmt.Println(f.Name())
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
