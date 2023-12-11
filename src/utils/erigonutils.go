package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// RequestBody represents the JSON payload for the request.
type RequestBody struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// Define the struct for each stage
type Stage struct {
	StageName   string `json:"stage_name"`
	BlockNumber string `json:"block_number"`
}

// Define the struct for the result field
type Result struct {
	CurrentBlock string  `json:"currentBlock"`
	HighestBlock string  `json:"highestBlock"`
	Stages       []Stage `json:"stages"`
}

// Define the overall struct for the JSON response
type ResponseBody struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Result `json:"result"`
}

type ResponseBodyNotSyncing struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  bool   `json:"result"`
}

type SyncingStatus struct {
	Progress float32
	Stage    string
}

// GetErigonSyncingStatus performs a request and returns the node status or an error.
func GetErigonSyncingStatus() (*SyncingStatus, error) {
	url := "http://localhost:8545/"
	requestBody := RequestBody{
		Jsonrpc: "2.0",
		Method:  "eth_syncing",
		Params:  []interface{}{},
		ID:      1,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	isSyncing := true
	var response ResponseBody
	if err := json.Unmarshal(respBody, &response); err != nil {
		// if already synced this is only returning false
		// so we can ignore the error
		var syncing ResponseBodyNotSyncing
		if err := json.Unmarshal(respBody, &syncing); err != nil {
			return nil, fmt.Errorf("error unmarshalling response body: %v", err)
		}
		if syncing.Result {
			// this is something else, report error
			return nil, fmt.Errorf("error unmarshalling response body: %v", err)
		}
		isSyncing = false
	}

	result := SyncingStatus{
		Progress: 100.0,
		Stage:    "Synced",
	}

	if isSyncing {
		highestBlock, _ := strconv.ParseInt(response.Result.HighestBlock, 16, 64)
		currentBlock, _ := strconv.ParseInt(response.Result.CurrentBlock, 16, 64)
		if highestBlock == 0 {
			result.Progress = 0
			result.Stage = "Waiting for Heimdall sync"

			// is it because fetching snapshots?
			logs, err := FetchContainerLogs("erigon", 10)
			if err != nil {
				return nil, fmt.Errorf("error fetching container logs: %v", err)
			}
			progress := findLastProgressUpdateInLogs(logs)
			if progress != nil {
				//result.Progress, _ = strconv.ParseFloat(progress["progress"], 32)
				result.Stage = fmt.Sprintf("Downloading snapshots [%s/%s]:%s", progress["step"], progress["total_steps"], progress["stage"])
			} else {
				// no its just synced apparently
			}
		} else {
			result.Progress = float32(currentBlock) / float32(highestBlock) * 100
			result.Stage = "Syncing"
		}
	}

	return &result, nil
}

type responseChainIDBody struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

// GetErigonChainID performs a request and returns the node status or an error.
func GetErigonChainID() (int, error) {
	url := "http://localhost:8545/"
	requestBody := RequestBody{
		Jsonrpc: "2.0",
		Method:  "eth_chainId",
		Params:  []interface{}{},
		ID:      1,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var response responseChainIDBody
	if err := json.Unmarshal(respBody, &response); err != nil {
		return 0, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	id, _ := strconv.ParseInt(strings.TrimPrefix(response.Result, "0x"), 16, 64)

	return int(id), nil
}

// findLastProgressUpdateInLogs searches for the last occurrence of a specific pattern in a string and extracts information.
func findLastProgressUpdateInLogs(input string) map[string]string {
	// Define the regular expression pattern to match both progress update formats
	pattern := regexp.MustCompile(`\[(\d+)/(\d+) ([^\]]+)\]( downloading\s+progress="([^"]+)" time-left=([^\s]+)| Waiting for torrents metadata: (\d+)/(\d+)| Indexing| Total)`)

	// Find all matches
	matches := pattern.FindAllStringSubmatch(input, -1)

	// Check if there are any matches
	if len(matches) == 0 {
		return nil
	}

	// Get the last match
	lastMatch := matches[len(matches)-1]

	// override the stage name if it's "Total" as its still part of the indexing process
	if lastMatch[4] == " Total" {
		lastMatch[4] = " Indexing"
	}

	// Check which format was matched and extract relevant information
	info := make(map[string]string)
	info["step"] = lastMatch[1]
	info["total_steps"] = lastMatch[2]
	info["stage"] = lastMatch[4]

	return info
}
