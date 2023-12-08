package utils

import (
	"encoding/json"
	"io"
	"net/http"
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
			LatestBlockHash   string `json:"latest_block_hash"`
			LatestAppHash     string `json:"latest_app_hash"`
			LatestBlockHeight string `json:"latest_block_height"`
			LatestBlockTime   string `json:"latest_block_time"`
			CatchingUp        bool   `json:"catching_up"`
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

// getNodeStatus performs an HTTP GET request to the specified URL and parses the JSON response.
func GetHeimdallNodeStatus() (*NodeStatusResponse, error) {
	resp, err := http.Get("http://localhost:26657/status")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return nil, err
	}

	var statusResponse NodeStatusResponse
	if err := json.Unmarshal(body, &statusResponse); err != nil {
		return nil, err
	}

	return &statusResponse, nil
}
