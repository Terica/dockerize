package impl

import (
	"io/ioutil"
	"os"
	"runtime"
)

// Runner interface for our testable implementations
type Runner interface {
	ReadFile(string) ([]byte, error)
	LookupEnv(string) (string, bool)
}

// RealRunner implements Runner (somehow)
type RealRunner struct{}

// ReadFile is an Implementation of ioutil.ReadFile
func (r RealRunner) ReadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

// LookupEnv is an Implementation of os.LookupEnv
func (r RealRunner) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

var (
	// GOOS is our abstracted OS version
	GOOS = runtime.GOOS
)
