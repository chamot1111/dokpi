package app

import "testing"

func TestPluginsCollection(t *testing.T) {
	pc, err := newPluginsCollection()
	if err != nil {
		t.Fatal(err)
	}
	if !pc.containsPlugin("test") {
		t.Fatal("couldn't find test plugin")
	}
}

func TestPluginDetect(t *testing.T) {
	p := plugin("test")
	d, err := p.detect(HelloRepo())
	if err != nil {
		t.Error("the buildpack return an error", err)
	}
	if !d {
		t.Error("the buildpack doesn't detect the unittest file")
	}
}
