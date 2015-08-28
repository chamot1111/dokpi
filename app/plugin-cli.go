package app

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	// PluginFolder is the folder where all the plugins are cloned.
	PluginFolder = "/home/git/.dokpi/plugins"
	// ExeInstall name of the file
	ExeInstall = "install"
	// ExeDetect name of the file
	ExeDetect = "detect"
	// ExeBuild name of the file
	ExeBuild = "build"
	// ExeDeploy name of the file
	ExeDeploy = "deploy"
)

// PluginExe describes the mandatory exe in a plugin repository
var PluginExe = [...]string{ExeInstall, ExeDetect, ExeBuild, ExeDeploy}

func init() {
	os.MkdirAll(PluginFolder, 644)
}

// PluginManager add plugins into the PluginFolder
// folder. The plugins are github repo.
//type PluginManager struct {
//}

// isGithubRepo check that the repo name has the pattern
// github.com/user/repo"
func isGithubRepo(repoName string) error {
	err := fmt.Errorf("The repo name %s should follow the pattern github.com/user/repo", repoName)
	if !strings.HasPrefix(repoName, "github.com/") {
		return err
	}
	components := strings.Split(repoName, "/")
	if len(components) != 3 {
		return err
	}
	return nil
}

func isLocalGitRepo(repoName string) bool {
	return FolderFileExists(repoName)
}

func ensurePluginFolder() error {
	if _, err := os.Stat(PluginFolder); os.IsNotExist(err) {
		err = os.Mkdir(PluginFolder, 644)
	}
	return nil
}

func url4Repo(repoName string) (url string, err error) {
	url = ""
	err = isGithubRepo(repoName)
	if err != nil {
		if !isLocalGitRepo(repoName) {
			err = fmt.Errorf("%s is not a local repository", repoName)
		} else {
			err = nil
			url = repoName
		}
	} else {
		err = nil
		url = "https://" + repoName + ".git"
	}
	return
}

func pluginName(repoName string) (name string, err error) {
	name = ""
	err = isGithubRepo(repoName)
	if err == nil {
		components := strings.Split(repoName, "/")
		name = components[2]
	} else if isLocalGitRepo(repoName) {
		err = nil
		name = path.Base(repoName)
		if name == "" {
			err = fmt.Errorf("can't clone from relative path")
		}
	} else {
		err = fmt.Errorf("%s repository path has error", repoName)
	}
	return
}

func pluginFolder(name string) string {
	return path.Join(PluginFolder, name)
}

func existsPlugin(name string) bool {
	_, err := os.Stat(pluginFolder(name))
	return os.IsExist(err)
}

func hasPluginError(name string) error {
	folder := pluginFolder(name)
	for _, exe := range PluginExe {
		if _, err := os.Stat(path.Join(folder, exe)); os.IsNotExist(err) {
			return fmt.Errorf("The file %s is lacking in the plugin %s", exe, name)
		}
	}
	return nil
}

func installPlugin(name string) error {
	installScriptPath := path.Join(pluginFolder(name), ExeInstall)
	if _, err := os.Stat(installScriptPath); os.IsNotExist(err) {
		return fmt.Errorf("The install file in the plugin %s doesn't exist", name)
	}

	cmd := exec.Command(installScriptPath)
	cmd.Dir = pluginFolder(name)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// AddPlugin clone the repo into the PluginFolder
func AddPlugin(repoName string) error {
	var err error
	err = ensurePluginFolder()
	if err != nil {
		return err
	}

	url, err := url4Repo(repoName)
	if err != nil {
		return err
	}

	name, err := pluginName(repoName)
	if err != nil {
		return err
	}

	if existsPlugin(name) {
		return fmt.Errorf("The plugin %s already exists", name)
	}

	destFolder := pluginFolder(name)
	err = os.Mkdir(destFolder, 644)
	if err != nil {
		return err
	}

	err = checkoutAction(url, destFolder)
	if err != nil {
		return err
	}

	err = installPlugin(name)
	if err != nil {
		return err
	}

	err = hasPluginError(name)
	if err != nil {
		return err
	}

	return nil
}
