package main

import (
  "syscall"
  "strings"
  "log"
  "os"
  "os/exec"
)

// Trampoline used to be sufficient until I needed environment handling

func main() {
        env := os.Environ()
        args := []string{}
	workdir := ""
	stage := 0
	for _, arg := range os.Args[1:] {
	  if stage == 0 && strings.Contains(arg, "=") {
		env = append(env, arg)
	  } else {
		if stage == 0 {
			workdir = arg
			stage++
		} else {
			args = append(args, arg)
		}
	  }
	}
	if len(args) == 0 {
		log.Fatal("Usage: execwdve [<env>=<val>]... <workdir> <cmd> [<args>]...")
	}
	err := syscall.Chdir(workdir)
	if err != nil {
		log.Fatalf("Can't change to working directory '%s'\n%s", workdir, err)
	}
	executable := ""
	if strings.Contains(args[0], "/") {
		executable = args[0]
	} else {
		executable, err = exec.LookPath(args[0])
		if err != nil {
			log.Fatalf("Can't find '%s' in the path\n%s", args[0], err)
		}
	}
	err = syscall.Exec(executable, args, env)
	log.Println(err)
}
