// Package docker is a convenience wrapper for fsouza's go-dockerclient
package docker

import (
	"github.com/fayep/dockerize/go/progress"
)

// It uses an anonymous docker.Client to store the real handle
// Client wraps a docker client
type MockClient struct {
	containers []APIContainers
	images     []APIContainers
	flags      map[string]bool
}

// Connect connects you to docker via the environment
func MockConnect() Mocker {
	return &MockClient{}
}

// Test that we have implemented the interface here at compile time
// This code will compile away to nothing.
// The linter won't error here even if the interface isn't fully implemented
var _ Mocker = (*MockClient)(nil)

// PStat gets you a list of running containers
func (cli *MockClient) PStat(filters map[string][]string) []APIContainers {
	return cli.containers
}

// Pull retrieves a container image from a repository
func (cli *MockClient) Pull(image string, tag string, pb *progress.Progress) {
}

// Run a container
// env represents additional environment variables
// mnts maps to binds because that's obvious.
func (cli *MockClient) Run(imageID string, name string, mnts []string, env []string, cmd []string) (string, error) {
	return imageID, nil
}

// Exec something in an existing container
func (cli *MockClient) Exec(container string, env []string, wd string, cmd []string) (int, error) {
	return 0, nil
}

func (cli *MockClient) AddContainer(cont APIContainers) {
	cli.containers = append(cli.containers, cont)
}

func (cli *MockClient) AddImage(cont APIContainers) {
	cli.images = append(cli.images, cont)
}

func (cli *MockClient) RemoveContainer(cont string) {
	containers := []APIContainers{}
	for _, c := range cli.containers {
		if c.ID != cont {
			containers = append(containers, c)
		}
	}
	cli.containers = containers
}

func (cli *MockClient) RemoveImage(cont string) {
	images := []APIContainers{}
	for _, i := range cli.images {
		if i.Image != cont {
			images = append(images, i)
		}
	}
	cli.images = images
}

func (cli *MockClient) ClearFlags() {
	cli.flags = map[string]bool{}
}

func (cli *MockClient) GetFlag(flag string) bool {
	return cli.flags[flag]
}

func (cli *MockClient) SetFlag(flag string, state bool) {
	cli.flags[flag] = state
}
