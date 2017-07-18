// Package docker for running stuff in containers
package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Container is a structure for storing container information
type Container struct {
	id string
	// We're computing the instance name
	Image   string
	Tag     string
	Volumes []string
}

// StandardName returns the standard name for the container
func (c Container) getName() (result string) {
	result = strings.Replace(c.Image, "/", "__", -1)
	result += "_" + c.getTag()
	return
}

// getTag returns a sane default for the tag or the specific one that has been set
func (c Container) getTag() (containerVersion string) {
	if c.Tag != "" {
		containerVersion = c.Tag
	} else {
		versionFile := fmt.Sprintf(".%s-version", c.Image)
		prefix := c.Image + "-"
		containerVersion = loadup(versionfile)
		strings.TrimPrefix(containerVersion, prefix)
	}
	if containerVersion == "" {
		containerVersion = "latest"
	}
	c.Tag = containerVersion
	// TODO(faye): this isn't the right place to set this, but I don't want to forget it.
	os.Setenv("container_version", containerVersion)
	return
}

// IsRunning determines if a given container is running by using docker ps
// By side effect, it may update the container id
func (c Container) IsRunning() (bool, error) {
	instanceName := c.getName()
	cmdOutput, err := exec.Command("docker", "ps", "-qf", "name="+instanceName).Output()
	if err != nil {
		return false, err
	}
	id := strings.TrimRight(string(cmdOutput), " \t\n\r")
	if id != "" {
		c.id = id
		return true, nil
	}
	return false, nil
}

// Run a container with some mounts
func (c Container) Run() (string, error) {
	running, err := c.IsRunning()
	if err != nil {
		return "", err
	}

	if running {
		return c.id, nil
	}
	containerVolumes := "$(pinata_mount_$os) $container_volumes"
	cleanup := "$(docker rm ${clean_container}_${container_version} 2>/dev/null)"
	instance := "$(docker run -td $container_volumes $environment --name ${clean_container}_${container_version} $container:$container_version cat)"
	c.id = instance
	return c.id, nil
}

// Exec a program in a container
func (c Container) Exec() {}
