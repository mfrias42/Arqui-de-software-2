package docker

import (
	"log"

	"github.com/docker/docker/client"
)

type DockerClient struct {
	Client *client.Client
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	log.Println("Docker client initialized successfully")
	return &DockerClient{Client: cli}, nil
}

func (d *DockerClient) Close() error {
	return d.Client.Close()
}
