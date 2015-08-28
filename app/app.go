package app

import (
	"fmt"
	"os"
	"os/exec"
)

// App run a git repo as a service.
type App struct {
	// Env for the app. If nil, the current environnement is use
	Env []string

	// GitRepoPath is the path of the git repo that contains the app
	GitRepoPath string

	// LogPath is the path of the log file
	LogPath string

	// StopChannel receive result of wait when the process finish
	StopChannel chan error

	ProcessState *os.ProcessState

	// folder in which the head of the git repo has been checkouted
	checkoutFolder string

	appPlugin plugin

	deployCmd *exec.Cmd
	running   bool
}

// NewApp create a new App object
func NewApp(gitRepoPath string) *App {
	return &App{
		GitRepoPath: gitRepoPath,
		StopChannel: make(chan error, 1),
	}
}

// checkout take a git repo path, and checkout it into a tmp folder.
func (a *App) checkout() error {
	var err error
	fmt.Println("checkout ", a.GitRepoPath)
	a.checkoutFolder, err = checkoutToTmp(a.GitRepoPath)
	return err
}

// Build trigger the build action of the folder. It depends of the technology.
func (a *App) build() error {
	return a.appPlugin.build(a.checkoutFolder)
}

// Deploy run the content of the folder depending of the content.
// It detects the technology in the repo, and run the corresponding command.
func (a *App) deploy() error {
	if a.deployCmd != nil {
		return fmt.Errorf("app %s is already running", a.GitRepoPath)
	}
	go a.waitEndProcess()
	return nil
}

func (a *App) waitEndProcess() {
	var err error
	var stdout *os.File
	stdout = nil

	if a.LogPath != "" {
		stdout, err := os.Create(a.LogPath)
		if err != nil {
			fmt.Println("Impossible to create log file", a.LogPath)
		}
		defer stdout.Close()
	}

	a.deployCmd, err = a.appPlugin.deploy(a.checkoutFolder, a.Env, stdout)
	if err != nil {
		a.StopChannel <- err
		return
	}

	a.running = true
	err = a.deployCmd.Wait()
	a.running = false
	a.ProcessState = a.deployCmd.ProcessState
	a.deployCmd = nil
	a.StopChannel <- err
}

func (a *App) retrieveAppPlugin() error {
	pluginsCollection, err := newPluginsCollection()
	if err != nil {
		return err
	}
	a.appPlugin = pluginsCollection.getPluginForRepo(a.checkoutFolder)
	return nil
}

// Run checkouts, build and deploy the git repo.
func (a *App) Run() error {
	var err error
	if a.IsRunning() {
		err = a.Kill()
		if err != nil {
			return err
		}
	}

	err = a.checkout()
	if err != nil {
		return err
	}

	err = a.retrieveAppPlugin()
	if err != nil {
		return err
	}
	if a.appPlugin.isNil() {
		return fmt.Errorf("couldn't find a plugin for the app %s", a.GitRepoPath)
	}

	err = a.build()
	if err != nil {
		return err
	}
	err = a.deploy()
	if err != nil {
		return err
	}
	return nil
}

// Kill the process of the app
func (a *App) Kill() error {
	if a.deployCmd != nil {
		return a.deployCmd.Process.Kill()
	}
	return nil
}

// IsRunning return true if the process is running
func (a *App) IsRunning() bool {
	return a.running
}
