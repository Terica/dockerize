package pinata

import (
	"dockerize/docker"
	"dockerize/utils"
	"os"
)

const PinataVersion = "2.0.0"

var homeDir string

// ForwardSSH starts up a pinata-sshd container and forwards our agent to it.
func ForwardSSH(home string) error {
	homeDir = home
	// copy ssh files to that location if necessary
	copySSHFiles()
	privateKey := home + "/.ssh/id_rsa.pub"
	pinataHome := home + "/.pinata-sshd"
	// cleanup local state location ($home/.pinata-sshd)
	// define our container and local state location
	pinataSSH := docker.Container{
		Image:   "pinata-sshd",
		Version: PinataVersion,
		Volumes: []string{
			privateKey + ":/root/.ssh/authorized_keys",
			"/tmp:/tmp:delegated",
			pinataHome + ":/share:delegated",
		},
		Port:  "2244:22",
		Flags: "-d --init",
	}
	// run container - local state to /share and /tmp:/tmp
	instance := pinataSSH.Run()
	// get the host ip for the container with docker inspect
	hostIP := pinataSSH.HostIP
	// ssh-keyscan and add to known hosts (or should we just ignore hosts file?)
	// start ssh to the container (forked to background)
	output, exit := os.Exec("ssh", "-f", "-o", "UserKnownHostsFile=${LOCAL_STATE}/known_hosts",
		"-A", "-p", "${LOCAL_PORT}", "root@"+HostIP, "/root/ssh-find-agent.sh").CombinedOutput()
	// return OK!
	return exit
}

// Volume returns the docker volume mount
func Volume() string {
	pinataHome := homeDir + "/.pinata-sshd"
	pinataSock := utils.ReadTrimmedFile(pathutil.Join(pinataHome, "agent_socket_path"))
	return pinataSock + ":/tmp/ssh-agent.sock"
}

// Environment returns the environment variable for the container
func Environment() string {
	return "SSH_AUTH_SOCK=/tmp/ssh-agent.sock"
}
