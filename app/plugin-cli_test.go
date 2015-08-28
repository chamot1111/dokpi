package app

import (
	"path"
	"testing"
)

func TestAddPluginBadRepo(t *testing.T) {
	rep := "test.com/toto/repo"
	err := AddPlugin(rep)
	if err == nil {
		t.Error("The repo ", rep, " should not be recognize")
	}
}

func TestPluginFolder(t *testing.T) {
	rep := "github.com/chamot1111/dokpi_go"
	name, err := pluginName(rep)
	if err != nil {
		t.Error(err)
	}
	computedPath := pluginFolder(name)
	dokpiFolder := path.Join(PluginFolder, "dokpi_go")
	if computedPath != dokpiFolder {
		t.Error("expected ", dokpiFolder, ", got ", computedPath)
	}
}

func TestAddPlugin(t *testing.T) {
	rep := "/tests-repo/hello"
	url, _ := url4Repo(rep)

	t.Log("add plugin", rep, "from", url)

	err := AddPlugin(rep)
	if err != nil {
		t.Error("Add plugin error: ", err)
	}
	if !folderFileExists(PluginFolder) {
		t.Error("the plugin folder doesn't exists ", PluginFolder)
	}
	dokpiFolder := path.Join(PluginFolder, "hello")
	if !folderFileExists(dokpiFolder) {
		t.Error("the plugin folder doesn't exists ", dokpiFolder)
	}
	installScriptPath := path.Join(dokpiFolder, ExeInstall)
	if !folderFileExists(installScriptPath) {
		t.Error("the install script doesn't exists ", installScriptPath)
	}
}
