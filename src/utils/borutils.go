package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// ResponseBody represents the expected JSON response.
type ResponseBody struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		CurrentBlock        string `json:"currentBlock"`
		HealedBytecodeBytes string `json:"healedBytecodeBytes"`
		HealedBytecodes     string `json:"healedBytecodes"`
		HealedTrienodeBytes string `json:"healedTrienodeBytes"`
		HealedTrienodes     string `json:"healedTrienodes"`
		HealingBytecode     string `json:"healingBytecode"`
		HealingTrienodes    string `json:"healingTrienodes"`
		HighestBlock        string `json:"highestBlock"`
		StartingBlock       string `json:"startingBlock"`
		SyncedAccountBytes  string `json:"syncedAccountBytes"`
		SyncedAccounts      string `json:"syncedAccounts"`
		SyncedBytecodeBytes string `json:"syncedBytecodeBytes"`
		SyncedBytecodes     string `json:"syncedBytecodes"`
		SyncedStorage       string `json:"syncedStorage"`
		SyncedStorageBytes  string `json:"syncedStorageBytes"`
	} `json:"result"`
	CatchingUp bool
}

// GetNodeStatus performs a request and returns the node status or an error.
func GetBorNodeStatus() (*ResponseBody, error) {
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

	var response ResponseBody
	response.CatchingUp = true
	if err := json.Unmarshal(respBody, &response); err != nil {
		// if already synced this is only returning false
		// so we can ignore the error
		var syncing struct {
			Jsonrpc string `json:"jsonrpc"`
			ID      int    `json:"id"`
			Result  bool   `json:"result"`
		}
		response.CatchingUp = false
		if err := json.Unmarshal(respBody, &syncing); err != nil {
			return nil, fmt.Errorf("error unmarshalling response body: %v", err)
		}
		if syncing.Result {
			// this is something else, report error
			return nil, fmt.Errorf("error unmarshalling response body: %v", err)
		}
	}

	return &response, nil
}

type responseChainIDBody struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

// GetNodeStatus performs a request and returns the node status or an error.
func GetBorChainID() (int, error) {
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
