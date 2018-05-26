// Package docker is a convenience wrapper for fsouza's go-dockerclient
package docker

import (
	"testing"
)

func TestMockConnect(t *testing.T) {
	cli := MockConnect()
	if nil == cli {
		t.Errorf("MockConnect returned nil")
	}
	if len(cli.(*MockClient).containers) != 0 {
		t.Errorf("There shouldn't be containers")
	}
}

func TestAddImage(t *testing.T) {
	cli := MockConnect()
	cli.AddImage(APIContainers{ID: NewID()})
	l := len(cli.(*MockClient).images)
	if l != 1 {
		t.Errorf("There should be 1 container, there are %d", l)
	}
}

func TestPStat(t *testing.T) {
	cli := MockConnect()
	cli.AddImage(APIContainers{ID: NewID()})
	result := cli.PStat(map[string][]string{})
	if len(result) != 1 {
		t.Errorf("len(result) should be 1: %d\n", len(result))
	}
}

func TestPull(t *testing.T) {
	cli := MockConnect()
	cli.Pull("container", "latest", nil)
	cont := cli.(*MockClient).containers
	if cont[0].Image != "container:latest" {
		t.Errorf("Not the container we are looking for: %s", cont[0].Image)
	}
}
