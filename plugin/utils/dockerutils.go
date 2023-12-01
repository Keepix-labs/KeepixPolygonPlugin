package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func CheckDockerExists() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false
	}
	defer cli.Close()

	_, err = cli.Info(context.Background())
	return err == nil
}

// DockerRun runs a Docker container and captures its output or if container is async, return container id.
func DockerRun(imageName string, args []string, containerPath, hostPath string, openPorts []uint, autorestart bool, networkName string, async bool, name string, sameUser bool) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	defer cli.Close()

	// Port bindings
	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for _, port := range openPorts {
		portStr := strconv.FormatUint(uint64(port), 10)
		portBindings[nat.Port(portStr+"/tcp")] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: portStr}}
		exposedPorts[nat.Port(portStr+"/tcp")] = struct{}{}
	}

	// Container configuration
	config := &container.Config{
		Image:        imageName,
		ExposedPorts: exposedPorts,
		Cmd:          args,
		Tty:          false,
	}

	if sameUser {
		config.User = fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid())
	}

	// Host configuration
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: hostPath,
				Target: containerPath,
			},
		},
	}
	if autorestart {
		hostConfig.RestartPolicy = container.RestartPolicy{Name: "unless-stopped"}
	}

	// Network configuration
	networkConfig := &network.NetworkingConfig{}
	if networkName != "" {
		networkConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			networkName: {},
		}
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, networkConfig, nil, name)
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	if async {
		return resp.ID, nil
	} else {
		// auto remove container after execution
		defer cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	}

	// Wait for the container to finish running
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: false})
	if err != nil {
		panic(err)
	}

	defer out.Close()
	content, err := io.ReadAll(io.Reader(out))
	if err != nil {
		panic(err)
	}

	return processDockerOutput(content), nil
}

func PullImage(imageName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, out)
	return err
}

// RemoveImage removes a Docker image and any containers created from it.
func RemoveImage(imageName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	// Find and remove containers that were created from the image
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if container.Image == imageName {
			if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true, RemoveVolumes: true}); err != nil {
				return fmt.Errorf("error removing container %s: %v", container.ID, err)
			}
		}
	}

	// Now remove the image
	_, err = cli.ImageRemove(ctx, imageName, types.ImageRemoveOptions{Force: true})
	if err != nil {
		return err
	}
	return nil
}

// CreateDockerNetwork creates a Docker network with the specified name.
func CreateDockerNetwork(networkName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	_, err = cli.NetworkCreate(ctx, networkName, types.NetworkCreate{
		CheckDuplicate: true, // Check for duplicate network names
	})
	if err != nil {
		return err
	}

	return nil
}

// RemoveDockerNetwork creates a Docker network with the specified name.
func RemoveDockerNetwork(networkName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	err = cli.NetworkRemove(ctx, networkName)
	if err != nil {
		return err
	}

	return nil
}

func processDockerOutput(log []byte) string {
	// Docker log output includes a header. For stdout, the first byte is 1.
	// We can use this information to strip the header.
	var cleanedLog []byte
	for i := 0; i < len(log); i++ {
		// Docker header size is 8 bytes. If the first byte is 1 (stdout), skip the header.
		if i+8 < len(log) && log[i] == 1 {
			cleanedLog = append(cleanedLog, log[i+8:]...)
			i += 7 // Skip the next 7 bytes of the header
		}
	}
	return string(cleanedLog)
}
