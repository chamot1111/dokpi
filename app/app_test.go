package app

import (
	"os"
	"path"
	"testing"
)

func HelloRepo() string {
	return path.Join(os.Getenv("TEST_REPO"), "/hello")
}

func TestCheckoutDetect(t *testing.T) {
	app := NewApp(HelloRepo())
	err := app.checkout()
	if err != nil {
		t.Fatal(err)
	}
	err = app.retrieveAppPlugin()
	if err != nil {
		t.Fatal(err)
	}
	if app.appPlugin.isNil() {
		t.Fatal("the plugin is empty")
	}
}

func TestRun(t *testing.T) {
	app := NewApp(HelloRepo())
	err := app.Run()
	if app.appPlugin.isNil() {
		t.Fatal("bad plugin")
	}
	if err != nil {
		t.Fatal(err)
	}
	err = <-app.StopChannel
	if err != nil {
		t.Fatal(err)
	}
}
