package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

const NDP_CONTAINER_NAME = "NGINX_DOCKER_PROXY"

var cli *client.Client

func initializeDockerClient() {
	cli = must(client.NewClientWithOpts(client.FromEnv))
}

func findNDPContainer() *types.Container {
	containers := must(cli.ContainerList(context.Background(), types.ContainerListOptions{All: true}))
	for i := 0; i < len(containers); i++ {
		if hasName(containers[i], NDP_CONTAINER_NAME) {
			return &containers[i]
		}
	}
	return nil
}

func createNDPContainer() {
	if isPortInUse(80) {
		panic("port 80 is in use")
	}
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "80",
	}
	formattedContainerPort := must(nat.NewPort("tcp", "80"))
	portBinding := nat.PortMap{formattedContainerPort: []nat.PortBinding{hostBinding}}
	createResponse := must(cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "nginx",
			ExposedPorts: nat.PortSet{formattedContainerPort: struct{}{}}},
		&container.HostConfig{PortBindings: portBinding},
		&network.NetworkingConfig{},
		&v1.Platform{},
		NDP_CONTAINER_NAME,
	))
	check(cli.ContainerStart(context.Background(), createResponse.ID, types.ContainerStartOptions{}))
}

func hasName(c types.Container, name string) bool {
	for i := 0; i < len(c.Names); i++ {
		if strings.Contains(c.Names[i], NDP_CONTAINER_NAME) {
			return true
		}
	}
	return false
}

func updateNginxConfInContainer(c *types.Container, newConf string) {
	file := must(os.Create(filepath.Join(os.TempDir(), "nginx.conf")))
	must(file.Write([]byte(newConf)))
	rc := must(archive.Tar(file.Name(), archive.Gzip))
	check(cli.CopyToContainer(context.Background(), c.ID, "/etc/nginx", rc, types.CopyToContainerOptions{}))
	file.Close()
	check(os.Remove(file.Name()))
}

func execNginxReloadInContainer(c *types.Container) {
	cmd := must(cli.ContainerExecCreate(context.Background(), c.ID, types.ExecConfig{
		Cmd: []string{"nginx", "-s", "reload"},
	}))
	check(cli.ContainerExecStart(context.Background(), cmd.ID, types.ExecStartCheck{}))
}

func isPortInUse(port int) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
