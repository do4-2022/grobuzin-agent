package main

import (
	"fmt"
	"os/exec"
)

func LaunchEngine(config Configuration, stopChan chan int) (err error) {

	var cmd *exec.Cmd

	switch config.Engine {
	case "nodejs":
		cmd = exec.Command("node", "/app/index.js")

	default:
		err = fmt.Errorf("engine %s not supported", config.Engine)

	}

	if err != nil {
		return
	}

	err = cmd.Start()

	if err != nil {
		return
	}

	<-stopChan

	err = cmd.Process.Kill()

	return
}
