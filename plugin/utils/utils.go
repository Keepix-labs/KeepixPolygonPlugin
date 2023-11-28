package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func CheckOSType() string {
	return runtime.GOOS
}

func CheckCPUCount() int {
	return runtime.NumCPU()
}

func CheckWSL2() bool {
	// WSL2 typically sets a WSL2 or WSL_DISTRO_NAME environment variable
	_, wsl2Exists := os.LookupEnv("WSL2")
	_, wslDistroExists := os.LookupEnv("WSL_DISTRO_NAME")

	return wsl2Exists || wslDistroExists
}

func CheckDockerExists() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false
	}
	defer cli.Close()

	_, err = cli.Info(context.Background())
	return err == nil
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

func RemoveImage(imageName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	defer cli.Close()

	_, err = cli.ImageRemove(ctx, imageName, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}
	return err
}

func WriteError(err string) {
	fmt.Fprintln(os.Stderr, err)
}
