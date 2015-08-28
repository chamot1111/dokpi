package app

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
)

func checkoutToTmp(gitRepoPath string) (destinationFolder string, err error) {
	destinationFolder, err = createTmpFolder(path.Base(gitRepoPath))
	if err != nil {
		return
	}
	err = checkoutAction(gitRepoPath, destinationFolder)
	if err != nil {
		return
	}
	return
}

func checkoutAction(repoFolder, destinationFolder string) error {
	out, err := exec.Command("git", "clone", repoFolder, destinationFolder).Output()
	if err != nil {
		return fmt.Errorf("%v: %v", err, string(out))
	}
	return nil
}

func createTmpFolder(repoName string) (path string, err error) {
	path, err = ioutil.TempDir("", "dokpi_"+repoName)
	return
}
