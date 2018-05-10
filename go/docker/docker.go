// +build !mocker

package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"bytes"
	"os"
	"io/ioutil"
	"local/dockerize/progress"
	"golang.org/x/crypto/ssh/terminal"

	docker "github.com/fsouza/go-dockerclient"
)

// ProgressDetail used by status updates of at least "Pull"
type ProgressDetail struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

// StatusUpdate used for docker status updates
type StatusUpdate struct {
	ID     string          `json:"id"`
	Status string          `json:"status"`
	Detail *ProgressDetail `json:"progressDetail"`
}

// Client wraps a docker client
type Client struct {
	*docker.Client
}

// Connect connects you to docker via the environment
func Connect() *Client {
	cli, _ := docker.NewClientFromEnv()
	return &Client{cli}
}

// PStat gets you a list of running containers
func (cli *Client) PStat(filters map[string][]string) []docker.APIContainers {
	listContainersOptions := docker.ListContainersOptions{
		Filters: filters,
		Context: context.Background(),
	}
	containers, _ := cli.ListContainers(listContainersOptions)
	return containers
}

// Pull retrieves a container image from a repository
func (cli *Client) Pull(image string, tag string, pb *progress.Progress) {
	status := new(bytes.Buffer)
	pullImageOptions := docker.PullImageOptions{
		Repository: image,
		Tag:	tag,
		OutputStream: status,
		RawJSONStream: true,
		Context: context.Background(),
	}
	authConfig := docker.AuthConfiguration{
	}
	if err := cli.PullImage(pullImageOptions,authConfig); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		manageProgress(status, pb)
	}
}

func manageProgress(status io.Reader, pb *progress.Progress) {
	if pb == nil {
		// Junk the status feed.  It has to be read until its end.
		io.Copy(ioutil.Discard, status)
	} else {
		decoder := json.NewDecoder(status)
		// We use Number to ensure that large numbers work ok.
		decoder.UseNumber()
		for decoder.More() {
			var (
				m StatusUpdate
			)
			decoder.Decode(&m)
			if m.Status == "Downloading" || m.Status == "Extracting" {
				if m.Status == "Downloading" {
					pb.OnlyAdd("Extracting "+m.ID, 0, m.Detail.Total)
				}
				pb.Add(m.Status+" "+m.ID, m.Detail.Current, m.Detail.Total)
				pb.Display()
			}
		}
		pb.Done()
	}
}

// Run a container
// env represents additional environment variables
// mnts maps to binds because that's obvious.
func (cli *Client) Run(imageID string, name string, mnts []string, env []string, cmd []string) (string, error){
	createContainerOptions := docker.CreateContainerOptions{
		Name: name,
		Context: context.Background(),
		Config: &docker.Config{
			Image: imageID,
			Cmd:   cmd,
			Tty:   true,
			Env:   env,
		},
		HostConfig: &docker.HostConfig{
			NetworkMode: "host",
			AutoRemove: false,
			Binds: mnts,
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
		AttachStdin: true,
		AttachStdout: true,
		AttachStderr: true,
		Tty: tty,
		Cmd: cmd,
		Env: env,
		WorkingDir: wd,
		Container: container,
		Context: context.Background(),
	}
	id := ""
	if resp, err := cli.CreateExec(createExecOptions); err != nil {
		return 255, err
	} else {
		id = resp.ID
	}
	startExecOptions := docker.StartExecOptions{
		OutputStream: os.Stdout,
		ErrorStream: os.Stderr,
		InputStream: os.Stdin,
		RawTerminal: false,
	}
	if err := cli.StartExec(id,startExecOptions); err != nil {
		return 255, err
	}
	if resp, err := cli.InspectExec(id); err != nil {
		return 255, err
	} else {
		//fmt.Printf("%+v\n", resp)
		return resp.ExitCode,nil
	}
}
