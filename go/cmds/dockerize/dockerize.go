package main

import (
	"fmt"
	"github.com/fayep/dockerize/go/docker"
	"os"
)

const dockerizeHelp = `
 ____             _             _
|  _ \  ___   ___| | _____ _ __(_)_______
| | | |/ _ \ / __| |/ / _ \ '__| |_  / _ \
| |_| | (_) | (__|   <  __/ |  | |/ /  __/
|____/ \___/ \___|_|\_\___|_|  |_/___\___|
Like Busybox for programming utilizing docker.Docker
docker.Docker and docker.Dockerize are my Desert Island Discs

With docker.Docker and docker.Dockerize you can do development within containers which
work as if they were part of your native development environment.

[download dockerize and docker on new laptop]
$ docker pull golang:latest
$ dockerize install golang go golint godoc gofmt
$ go help
Go is a tool for managing Go source code.

Usage:

        go command [arguments]
...

$ dockerize install ruby bundle irb
$ echo "2.3.1" > .ruby-version
$ bundle
Unable to find image 'ruby:2.3.1' locally
2.3.1: Pulling from library/ruby
Status: Downloaded newer image for ruby:2.3.1
Fetching gem metadata from https://rubygems.org/...........
Fetching version metadata from https://rubygems.org/..
Resolving dependencies...
...

Brought to you by the MPL 2.0

Author: Faye Salwin (beiriannydd) faye@futureadvisor.com
`

func init() {
	// Each dockerize module should add the "modes" aka dockerize <foo> that it supports.
	// That's how we can make sure that stuff only intended for a particular OS can only work on that OS.
	// If you +build !linux then the init() for that .go source won't be run either.
	programMode["dockerize"] = modeDockerize
}

func dockerizeInstall() int {
	return 0
}

func modeDockerize(cli docker.Docker) int {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			fmt.Println(dockerizeHelp)
		case "install":
			dockerizeInstall()
		}
	} else {
		fmt.Println(dockerizeHelp)
	}
	return 0
}
