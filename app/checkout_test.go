package app

import (
	"os"
	"path"
	"strings"
	"testing"
)

func TestCreateTmpFolder(t *testing.T) {
	folderPath, err := createTmpFolder("myrepo")
	if err != nil {
		t.Fatal("can't create tmp folder: ", err)
	}
	expectedName := "dokpi_myrepo"
	if !strings.HasPrefix(path.Base(folderPath), expectedName) {
		t.Error("expected prefix is ", expectedName, " got: ", folderPath)
	}
	if !FolderFileExists(folderPath) {
		t.Error("The folder has not been created")
	}
}

func TestCheckoutToTmp(t *testing.T) {
	destinationFolder, err := checkoutToTmp(os.Getenv("TEST_REPO") + "/hello")
	if err != nil {
		t.Error(err)
	}
	if !FolderFileExists(destinationFolder + "/unittest") {
		t.Error("git repo has not been cloned")
	}
}
