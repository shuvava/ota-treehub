package fshelper

import (
	"fmt"
	"os"
	"syscall"
)

func isPathExist(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path doesn't exist")
	}

	return info, nil
}

func isDirExist(path string) (os.FileInfo, error) {
	info, err := isPathExist(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path isn't a directory")
	}
	return info, nil
}

// IsPathExist checks if path exist
func IsPathExist(path string) bool {
	_, err := isPathExist(path)
	return err == nil
}

// IsDirExist checks if path exist and it is a directory
func IsDirExist(path string) error {
	_, err := isDirExist(path)
	return err
}

// IsDirWritable checks write permission on directory
func IsDirWritable(path string) (bool, error) {
	info, err := isDirExist(path)
	if err != nil {
		return false, err
	}
	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, fmt.Errorf("write permission bit is not set on this file for user")
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		return false, fmt.Errorf("unable to get stat")
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, fmt.Errorf("user doesn't have permission to write to this directory")
	}

	return true, nil
}

// EnsureDir creates directory or throw an error
func EnsureDir(path string) error {
	err := os.Mkdir(path, 0740)
	if err == nil {
		return nil
	}

	return IsDirExist(path)
}
