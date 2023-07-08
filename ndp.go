package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	defer mainRecover()
	if len(os.Args) != 2 {
		panic("usage: ndp <port>")
	}
	port := os.Args[1]
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		panic("usage: ndp <port>")
	}
	if portInt < 1024 || 65535 < portInt {
		panic("port must be between 1024-65535")
	}
	initializeDockerClient()
	c := findNDPContainer()
	if c == nil {
		createNDPContainer()
		c = findNDPContainer()
	}
	newConf := formatNginxConf(port)
	updateNginxConfInContainer(c, newConf)
	execNginxReloadInContainer(c)
}

func mainRecover() {
	if r := recover(); r != nil {
		fmt.Println(r)
		os.Exit(1)
	}
}

