package main

import (
	"fmt"
	"log"
	"os/exec"
)

func LaunchEngine(config Configuration, stopChan chan int) (err error) {

	var cmd *exec.Cmd

	switch config.Engine {
	case "nodejs":
		cmd = exec.Command("node", "index.js")
		cmd.Dir = "/app"

	default:
		err = fmt.Errorf("engine %s not supported", config.Engine)

	}

	if err != nil {
		log.Println(err)
		return
	}

	err = cmd.Start()

	if err != nil {
		log.Println("Error starting engine", err)
		return
	}

	<-stopChan

	err = cmd.Process.Kill()

	return
}
