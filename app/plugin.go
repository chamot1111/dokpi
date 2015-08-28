package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

// Plugin ...
type plugin string

// PluginsCollection ...
type pluginsCollection []plugin

// NewPluginsCollection inialize a PluginsCollection
func newPluginsCollection() (pc pluginsCollection, err error) {
	if !FolderFileExists(PluginFolder) {
		err = fmt.Errorf("The PluginFolder donesn't exists (%s)", PluginFolder)
		return
	}
	folders, err := ioutil.ReadDir(PluginFolder)
	if err != nil {
		return
	}
	pc = make([]plugin, 0)
	for _, f := range folders {
		if f.IsDir() {
			err = hasPluginError(f.Name())
			if err != nil {
				return
			}
			pc = append(pc, plugin(f.Name()))
		}
	}
	return
}

func (pc pluginsCollection) getPluginForRepo(repoPath string) plugin {
	for _, p := range pc {
		d, err := p.detect(repoPath)
		if err != nil {
			fmt.Println(err)
		} else if d {
			return p
		}
	}
	return ""
}

func (pc pluginsCollection) containsPlugin(name string) bool {
	for _, s := range pc {
		if name == string(s) {
			return true
		}
	}
	return false
}

func (p plugin) detect(repoPath string) (r bool, err error) {
	detectScript := path.Join(pluginFolder(string(p)), ExeDetect)
	detectCmd := exec.Command(detectScript)
	detectCmd.Dir = repoPath
	errcmd := detectCmd.Run()
	if errcmd == nil {
		return true, nil
	}
	if _, ok := errcmd.(*exec.ExitError); ok {
		return false, nil
	}
	return false, errcmd
}

func (p plugin) isNil() bool {
	return string(p) == ""
}

func (p plugin) build(repoPath string) error {
	buildScript := path.Join(pluginFolder(string(p)), ExeBuild)
	buildCmd := exec.Command(buildScript)
	buildCmd.Dir = repoPath
	return buildCmd.Run()
}

func (p plugin) deploy(repoPath string, env []string, stdout *os.File) (cmd *exec.Cmd, err error) {
	deployScript := path.Join(pluginFolder(string(p)), ExeDeploy)
	cmd = exec.Command(deployScript)
	cmd.Dir = repoPath
	cmd.Env = env
	cmd.Stdout = stdout
	err = cmd.Start()
	return
}
