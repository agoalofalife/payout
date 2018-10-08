package utils

import (
	"log"
	"os"
	"os/exec"
)

// check exist command or installed
func ExistCliCommand(command string) string {
	path, err := exec.LookPath(command)
	if err != nil {
		log.Fatalf("You need to install %s!", command)
		os.Exit(-1)
	}
	return path
}

// check exist file, if is not exist , output error
func ExistFile(pathToFile string, err error) {
	if _, isFile := os.Stat(pathToFile); os.IsNotExist(isFile) {
		panic(err.Error())
	}
}
