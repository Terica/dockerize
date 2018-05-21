package docker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fayep/dockerize/go/progress"
	"golang.org/x/crypto/ssh/terminal"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

// It uses an anonymous docker.Client to store the real handle
// Client wraps a docker client
type Client struct {
	*docker.Client
}

// Connect connects you to docker via the environment
func Connect() Docker {
	cli, _ := docker.NewClientFromEnv()
	return &Client{cli}
}

// Test that we have implemented the interface here
// This code will compile away to nothing.
var _ Docker = (*Client)(nil)

// PStat gets you a list of running containers
func (cli *Client) PStat(filters map[string][]string) []APIContainers {
	listContainersOptions := docker.ListContainersOptions{
		Filters: filters,
		Context: context.Background(),
	}
	cont, _ := cli.ListContainers(listContainersOptions)
	var containers []APIContainers
	DeepCopy(cont, containers)
	return containers
}

// Pull retrieves a container image from a repository
func (cli *Client) Pull(image string, tag string, pb *progress.Progress) {
	status := new(bytes.Buffer)
	pullImageOptions := docker.PullImageOptions{
		Repository:    image,
		Tag:           tag,
		OutputStream:  status,
		RawJSONStream: true,
		Context:       context.Background(),
	}
	authConfig := docker.AuthConfiguration{}
	if err := cli.PullImage(pullImageOptions, authConfig); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		manageProgress(status, pb)
	}
}

// Run a container
// env represents additional environment variables
// mnts maps to binds because that's obvious.
func (cli *Client) Run(imageID string, name string, mnts []string, env []string, cmd []string) (string, error) {
	createContainerOptions := docker.CreateContainerOptions{
		Name:    name,
		Context: context.Background(),
		Config: &docker.Config{
			Image: imageID,
			Cmd:   cmd,
			Tty:   true,
			Env:   env,
		},
		HostConfig: &docker.HostConfig{
			NetworkMode: "host",
			AutoRemove:  false,
			Binds:       mnts,
		},
	}
	resp, err := cli.CreateContainer(createContainerOptions)
	if err != nil {
		return "", err
	}
	err = cli.StartContainer(resp.ID, nil)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// Exec something in an existing container
func (cli *Client) Exec(container string, env []string, wd string, cmd []string) (int, error) {
	fd := int(os.Stdin.Fd())
	tty := false
	if terminal.IsTerminal(fd) {
		oldState, err := terminal.MakeRaw(fd)
		if err != nil {
			// handle err ...
		}
		defer terminal.Restore(fd, oldState)
		tty = true
	}
	createExecOptions := docker.CreateExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          tty,
		Cmd:          cmd,
		Env:          env,
		WorkingDir:   wd,
		Container:    container,
		Context:      context.Background(),
	}
	id := ""
	if resp, err := cli.CreateExec(createExecOptions); err != nil {
		return 255, err
	} else {
		id = resp.ID
	}
	startExecOptions := docker.StartExecOptions{
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		InputStream:  os.Stdin,
		RawTerminal:  false,
	}
	if err := cli.StartExec(id, startExecOptions); err != nil {
		return 255, err
	}
	if resp, err := cli.InspectExec(id); err != nil {
		return 255, err
	} else {
		//fmt.Printf("%+v\n", resp)
		return resp.ExitCode, nil
	}
}
