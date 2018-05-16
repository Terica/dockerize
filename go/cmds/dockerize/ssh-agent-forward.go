// +build linux

package main

import (
	"fmt"
	"github.com/fayep/dockerize/go/docker"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var forwardPath string

func init() {
	programMode["ssh-agent-forward"] = sshAgentForwarder
}

func sshAgentForwarder(cli docker.Docker) int {
	forwardPath = filepath.FromSlash("/share/agent_socket_path")
	if authSockPath, found := os.LookupEnv("SSH_AUTH_SOCK"); found {
		os.Chown(authSockPath, 1000, 1000)
		handleInterrupt(2)
		if err := ioutil.WriteFile(forwardPath, []byte(authSockPath+"\n"), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "writefile: %s\n", err.Error())
		}
		fmt.Printf("Found SSH-Agent-Socket at %s\n", authSockPath)
		fmt.Printf("Storing that location at %s until exit\n", forwardPath)
		// Discard any input - basically wait until interrupted
		io.Copy(ioutil.Discard, os.Stdin)
	}
	return 0
}

func handleInterrupt(intrptChSize int) {
	// Notify the channel s about signals
	s := make(chan os.Signal, intrptChSize)
	signal.Notify(s,
		syscall.SIGABRT, syscall.SIGBUS, syscall.SIGFPE,
		syscall.SIGHUP, syscall.SIGILL, syscall.SIGINT,
		syscall.SIGKILL, syscall.SIGPIPE, syscall.SIGPROF,
		syscall.SIGQUIT, syscall.SIGSEGV, syscall.SIGSTOP,
		syscall.SIGSYS, syscall.SIGTERM, syscall.SIGTRAP,
		syscall.SIGURG, syscall.SIGUSR1, syscall.SIGUSR2)
	// Process signal
	go func() {
		for sig := range s {
			fmt.Printf("exiting on %s signal\n", sig)
			//cleanup()
			os.Remove(forwardPath)
			os.Exit(0)
		}
	}()
}
