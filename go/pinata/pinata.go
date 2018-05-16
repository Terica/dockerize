package pinata

import (
	"fmt"
	"local/dockerize/docker"
	"os"
)

type pinata interface {
	Run()
}

func checkImage() {

}

func build() {
	fmt.Println("Please checkout and build https://github.com/FutureAdvisor/pinata-ssh-agent.git")
	os.Exit(1)
}

func checkRunning() bool {
}

// Run is the main entrypoint to pinata, we just want to make sure Pinata container is running
func Run() bool {
	if checkRunning() {
		return true
	}

	if checkImage() {
		docker.Run("pinata")
	} else {
		build()
	}
}
