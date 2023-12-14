package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// Define struct to match the JSON structure
type NodeStatusResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		NodeInfo struct {
			ProtocolVersion struct {
				P2P   string `json:"p2p"`
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"protocol_version"`
			ID         string            `json:"id"`
			ListenAddr string            `json:"listen_addr"`
			Network    string            `json:"network"`
			Version    string            `json:"version"`
			Channels   string            `json:"channels"`
			Moniker    string            `json:"moniker"`
			Other      map[string]string `json:"other"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHash    string `json:"latest_block_hash"`
			LatestAppHash      string `json:"latest_app_hash"`
			LatestBlockHeight  string `json:"latest_block_height"`
			CurrentBlockHeight string
			LatestBlockTime    string `json:"latest_block_time"`
			CatchingUp         bool   `json:"catching_up"`
		} `json:"sync_info"`
		ValidatorInfo struct {
			Address string `json:"address"`
			PubKey  struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
			VotingPower string `json:"voting_power"`
		} `json:"validator_info"`
	} `json:"result"`
}

// RunSnapshotDownloader starts the process to download the proper snapshot for heimdall
func RunSnapshotDownloader(hostHeimdallPath string, network string) error {
	err := PullImage("alpine:latest")
	if err != nil {
		return err
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating Docker client: %v", err)
	}
	defer cli.Close()

	command := fmt.Sprintf(`apk add aria2 curl bash zstd pv && curl -L https://snapshot-download.polygon.technology/snapdown.sh | bash -s -- --network %s --client heimdall --extract-dir /heimdall/data`, network)
	tempContainerConfig := container.Config{
		Image: "alpine:latest",
		Cmd:   []string{"sh", "-c", command},
	}

	hostConfig := container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: hostHeimdallPath,
				Target: "/heimdall",
			}},
	}

	resp, err := cli.ContainerCreate(ctx, &tempContainerConfig, &hostConfig, nil, nil, "heimdall-snapshot-downloader")
	if err != nil {
		return fmt.Errorf("error creating snapshot downloader container: %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("error starting snapshot downloader container: %v", err)
	}

	return nil
}

func ValidateSnapshot(hostHeimdallPath string) (bool, error) {
	running, err := IsContainerRunning("heimdall-snapshot-downloader")
	if err != nil {
		return false, fmt.Errorf("error checking status of snapshot downloader: %v", err)
	}
	if running {
		return false, nil
	}
	// get container logs
	_, err = FetchContainerLogs("heimdall-snapshot-downloader", 1)
	if err != nil {
		return false, fmt.Errorf("error getting logs from snapshot downloader: %v", err)
	}
	// check if last line of the output is a valid download

	// remove image and container
	err = RemoveImageIfExists("alpine:latest")
	if err != nil {
		return false, fmt.Errorf("error removing image: %v", err)
	}
	return true, nil
}

// extractProgressValue extracts the numeric value from a progress string like "(34%)"
func extractProgressValue(progress string) (int, error) {
	// Remove non-numeric characters
	re := regexp.MustCompile(`[^0-9]`)
	numericString := re.ReplaceAllString(progress, "")

	// Convert to integer
	return strconv.Atoi(numericString)
}

// getProgressFromProgressSummary searches for the last occurrence of the download summary pattern in Docker logs and extracts information.
func getProgressFromProgressSummary(input string) (float32, error) {
	// Define the regular expression pattern to match the download summary and file details
	// The pattern considers the binary prefix at the start of each line
	pattern := regexp.MustCompile(`\*\*\* Download Progress Summary as of .*? \*\*\* \n([\s\S]*?)\n\n`)

	// Normalize line endings and remove binary prefixes
	normalizedInput := strings.ReplaceAll(input, "\n\x01\x00\x00\x00\x00\x00\x00P", "\n")
	normalizedInput = strings.ReplaceAll(normalizedInput, "\n\x01\x00\x00\x00\x00\x00\x00F", "\n")
	normalizedInput = strings.ReplaceAll(normalizedInput, "\n\x01\x00\x00\x00\x00\x00\x00B", "\n")
	normalizedInput = strings.ReplaceAll(normalizedInput, "\n\x01\x00\x00\x00\x00\x00\x00C", "\n")
	normalizedInput = strings.ReplaceAll(normalizedInput, "\n\x01\x00\x00\x00\x00\x00\x00\x01", "\n")
	// Find all matches
	matches := pattern.FindAllStringSubmatch(normalizedInput, -1)

	// Check if there are any matches
	if len(matches) == 0 {
		return 0, fmt.Errorf("no matching download summary found")
	}

	// Get the last match
	lastMatch := matches[len(matches)-1]

	progressReport := lastMatch[1]

	patternProgressPercent := regexp.MustCompile(`\(\d+%\)`)
	matches = patternProgressPercent.FindAllStringSubmatch(progressReport, -1)

	totalProgress := 0

	for _, match := range matches {
		progressValue, err := extractProgressValue(match[0])
		if err != nil {
			return 0, fmt.Errorf("error extracting progress value: %v", err)
		}
		totalProgress += progressValue
	}

	totalFiles := len(matches)
	progress := float32(totalProgress) / float32(totalFiles*100) * 100
	return progress, nil
}

func SnapshotProgress() (float32, error) {
	logs, err := FetchContainerLogs("heimdall-snapshot-downloader", 100)
	if err != nil {
		return 0, fmt.Errorf("error getting logs from snapshot downloader: %v", err)
	}
	progress, err := getProgressFromProgressSummary(logs)
	if err != nil {
		return 0, fmt.Errorf("error getting progress from progress summary: %v", err)
	}
	return progress, nil
}

// getNodeStatus performs an HTTP GET request to the specified URL and parses the JSON response.
func GetHeimdallNodeStatus() (*NodeStatusResponse, error) {
	resp, err := http.Get("http://localhost:26657/status")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyStatus, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return nil, err
	}

	var statusResponse NodeStatusResponse
	if err := json.Unmarshal(bodyStatus, &statusResponse); err != nil {
		return nil, err
	}

	// also fetch live status to get the current block height
	// testnet
	rpc := "https://heimdall-api.polygon.technology/staking/validator-set"
	if statusResponse.Result.NodeInfo.Network != "heimdall-137" {
		// check testnet instead
		rpc = "https://heimdall-api-testnet.polygon.technology/staking/validator-set"
	}
	resp, err = http.Get(rpc)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyLiveStatus, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return nil, err
	}

	var liveStatusResponse = struct {
		Height string `json:"height"`
	}{}
	if err := json.Unmarshal(bodyLiveStatus, &liveStatusResponse); err != nil {
		return nil, err
	}

	statusResponse.Result.SyncInfo.CurrentBlockHeight = liveStatusResponse.Height

	return &statusResponse, nil
}
