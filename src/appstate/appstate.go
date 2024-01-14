package appstate

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

// AppStateEnum is a list of possible application states.
type AppStateEnum int

const (
	NoState AppStateEnum = iota
	SetupErrorState
	StartingInstall
	InstallingNode
	ConfiguringHeimdall
	ConfiguringErigon
	ConfiguringNetwork
	NodeInstalled
	StartingNode
	StartingHeimdall
	StartingRestServer
	StartingErigon
	NodeStarted
	NodeRestarting
	// Add new states here...
)

type Account struct {
	Address string `json:"address"`
	PK      string `json:"pk"`
}

type AppState struct {
	State                      AppStateEnum `json:"state"`
	IsTestnet                  bool         `json:"isTestnet"`
	HeimdallSnapshotDownloaded bool         `json:"heimdallSnapshotDownloaded"`
	Wallet                     Account      `json:"wallet"`
	RPC                        string       `json:"rpc"`
}

// CurrentState holds the current state of the application.
var CurrentState AppState = AppState{State: NoState, IsTestnet: false, Wallet: Account{Address: "", PK: ""}}

func CurrentStateString() string {
	switch CurrentState.State {
	case NoState:
		return "NoState"
	case SetupErrorState:
		return "SetupErrorState"
	case StartingInstall:
		return "StartingInstall"
	case InstallingNode:
		return "InstallingNode"
	case ConfiguringHeimdall:
		return "ConfiguringHeimdall"
	case ConfiguringErigon:
		return "ConfiguringErigon"
	case ConfiguringNetwork:
		return "ConfiguringNetwork"
	case NodeInstalled:
		return "NodeInstalled"
	case StartingNode:
		return "StartingNode"
	case StartingHeimdall:
		return "StartingHeimdall"
	case StartingRestServer:
		return "StartingRestServer"
	case StartingErigon:
		return "StartingErigon"
	case NodeStarted:
		return "NodeStarted"
	case NodeRestarting:
		return "NodeRestarting"
	}

	return "Unknown"
}

// UpdateState updates the current state and writes it to disk.
func UpdateState(newState AppStateEnum) error {
	CurrentState.State = newState
	return writeStateToFile(CurrentState)
}

// UpdateChain updates the current state and writes it to disk.
func UpdateChain(isTestnet bool) error {
	CurrentState.IsTestnet = isTestnet
	return writeStateToFile(CurrentState)
}

// UpdateAccount updates the current state and writes it to disk.
func UpdateAccount(privateKey string, publicKey string) error {
	CurrentState.Wallet.Address = publicKey
	CurrentState.Wallet.PK = privateKey
	return writeStateToFile(CurrentState)
}

// UpdateSnapshotDownloaded updates the current state and writes it to disk.
func UpdateSnapshotDownloaded(downloaded bool) error {
	CurrentState.HeimdallSnapshotDownloaded = downloaded
	return writeStateToFile(CurrentState)
}

// UpdateSnapshotDownloaded updates the current state and writes it to disk.
func UpdateRPC(rpc string) error {
	CurrentState.RPC = rpc
	return writeStateToFile(CurrentState)
}

// LoadState loads the current state from the file, if it exists.
func LoadState() error {
	path, err := GetStoragePath()
	if err != nil {
		return err
	}

	filePath := filepath.Join(path, "state.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// State file does not exist, no state to load
		return nil
	}

	stateJSON, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var state AppState
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		return err
	}

	CurrentState = state
	return nil
}

// writeStateToFile writes the current state to a file in JSON format.
func writeStateToFile(state AppState) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return err
	}

	path, err := GetStoragePath()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(path, "state.json"), stateJSON, fs.FileMode(0644))
}

// GetStoragePath gets the path to the storage directory.
func GetStoragePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	pluginFolder := filepath.Join(home, ".keepix/plugins/keepix-polygon-plugin", "/")
	err = os.MkdirAll(pluginFolder, os.ModePerm)
	if err != nil {
		return "", err
	}

	return pluginFolder, nil
}
