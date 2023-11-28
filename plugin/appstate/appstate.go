package appstate

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

// AppStateEnum is a list of possible application states.
type AppStateEnum string

const (
	NoState         AppStateEnum = "NO_STATE"
	StartingInstall AppStateEnum = "STARTING_INSTALLATION"
	InstallingCLI   AppStateEnum = "INSTALLING_CLI"
	InstallingNode  AppStateEnum = "INSTALLING_NODE"
	ConfiguringNode AppStateEnum = "CONFIGURING_NODE"
	StartingNode    AppStateEnum = "STARTING_NODE"
	NodeRunning     AppStateEnum = "NODE_RUNNING"
	NodeStopped     AppStateEnum = "NODE_STOPPED"
	NodeRestarting  AppStateEnum = "NODE_RESTARTING"
	SetupErrorState AppStateEnum = "SETUP_ERROR_STATE"
	// Add new states here...
)

// CurrentState holds the current state of the application.
var CurrentState AppStateEnum = NoState

// UpdateState updates the current state and writes it to disk.
func UpdateState(newState AppStateEnum) error {
	CurrentState = newState
	return writeStateToFile(CurrentState)
}

// LoadState loads the current state from the file, if it exists.
func LoadState() error {
	path, err := getStoragePath()
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

	var state AppStateEnum
	err = json.Unmarshal(stateJSON, &state)
	if err != nil {
		return err
	}

	CurrentState = state
	return nil
}

// writeStateToFile writes the current state to a file in JSON format.
func writeStateToFile(state AppStateEnum) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return err
	}

	path, err := getStoragePath()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(path, "state.json"), stateJSON, fs.FileMode(0644))
}

// getStoragePath gets the path to the storage directory.
func getStoragePath() (string, error) {
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
