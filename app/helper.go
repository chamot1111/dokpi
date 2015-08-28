package app

import "os"

func FolderFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
