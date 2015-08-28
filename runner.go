package dokpi

import (
	"fmt"
	"io/ioutil"

	"github.com/chamot1111/dokpi/app"
)

const repoFolder = "/home/git/.dokpi/repository"

type appsRunner struct {
	availableApps []string
	runningApps   map[string]*app.App
}

func newAppsRunner() (ar *appsRunner, err error) {
	ar = &appsRunner{
		availableApps: make([]string, 0),
		runningApps:   make(map[string]*app.App)}
	err = ar.updateAvailableApps()
	return
}

func (ar *appsRunner) removeDeadApps() {
	var delKeys []string
	for k, app := range ar.runningApps {
		if !app.IsRunning() {
			delKeys = append(delKeys, k)
		}
	}
	for _, k := range delKeys {
		delete(ar.runningApps, k)
	}
}

func (ar *appsRunner) startApp(name string) error {
	ar.removeDeadApps()
	ar.updateAvailableApps()
	if !ar.existsApp(name) {
		return fmt.Errorf("the app %s doesn't exist on the server", name)
	}
	if _, ok := ar.runningApps[name]; ok {
		return fmt.Errorf("the app %s is already running", name)
	}

	newApp := app.NewApp(name)
	err := newApp.Run()
	if err != nil {
		return err
	}

	ar.runningApps[name] = newApp
	return nil
}

func (ar *appsRunner) stopApp(name string) error {
	ar.removeDeadApps()
	app, exists := ar.runningApps[name]
	if !exists {
		return fmt.Errorf("the app %s is not running", name)
	}

	err := app.Kill()
	return err
}

func (ar *appsRunner) updateAvailableApps() error {
	if !app.FolderFileExists(repoFolder) {
		return fmt.Errorf("The repository folder doesn't exists (%s)", repoFolder)
	}
	folders, err := ioutil.ReadDir(repoFolder)
	if err != nil {
		return err
	}
	ar.availableApps = make([]string, 0)
	for _, f := range folders {
		if f.IsDir() {
			ar.availableApps = append(ar.availableApps, f.Name())
		}
	}
	return nil
}

func (ar *appsRunner) getAppsStatus() []appStatus {
	ret := make([]appStatus, 0, len(ar.availableApps))
	ar.removeDeadApps()
	ar.updateAvailableApps()
	for _, appName := range ar.availableApps {
		_, running := ar.runningApps[appName]
		ret = append(ret, appStatus{name: appName, running: running})
	}
	return ret
}

func (ar appsRunner) existsApp(name string) bool {
	for _, s := range ar.availableApps {
		if s == name {
			return true
		}
	}
	return false
}
