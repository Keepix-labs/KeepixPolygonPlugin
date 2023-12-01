package tasks

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/utils"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
)

// TaskFunc defines the signature for task functions
type TaskFunc func([]string) bool

// TaskMap maps task names to their corresponding functions
var TaskMap = map[string]TaskFunc{
	"install":    installTask,
	"uninstall":  uninstallTask,
	"installed":  installedTask,
	"status":     statusTask,
	"start":      startTask,
	"stop":       stopTask,
	"sync-state": syncStateTask,
}

// TaskRequirements maps task names to their required system conditions
var TaskRequirements = map[string][]string{
	"install":    {"docker", "uninstalled", "linux", "cpu4"},
	"uninstall":  {"docker", "installed", "stopped"},
	"installed":  {"docker"},
	"status":     {},
	"start":      {"docker", "stopped"},
	"stop":       {"docker", "running"},
	"sync-state": {"docker", "running"},
}

//go:embed conf/bor/genesis.json
var genesis string

//go:embed conf/bor/config.toml
var configToml string

// installTask is an example task for installation purposes
func installTask(args []string) bool {
	if len(args) != 1 {
		utils.WriteError("Invalid args")
		return false
	}

	ethereumRPC := args[0]
	if !utils.IsValidURL(ethereumRPC) {
		utils.WriteError("Invalid ethereumRPC")
		return false
	}

	storage, _ := appstate.GetStoragePath()
	localPathHeimdall := path.Join(storage, "data", "heimdall")
	localPathBor := path.Join(storage, "data", "bor")

	if appstate.CurrentState <= appstate.InstallingNode {
		// not installed yet

		appstate.UpdateState(appstate.InstallingNode)

		err := utils.PullImage("0xpolygon/heimdall:1.0.3")
		if err != nil {
			utils.WriteError("Error pulling heimdall image:" + err.Error())
			return false
		}

		err = utils.PullImage("0xpolygon/bor:1.1.0")
		if err != nil {
			utils.WriteError("Error pulling bor image:" + err.Error())
			return false
		}

		// setting up local config path
		err = os.MkdirAll(localPathHeimdall, os.ModePerm)
		if err != nil {
			utils.WriteError("Error creating local path:" + err.Error())
			return false
		}

		// check heimdall
		output, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"heimdallcli", "version"}, "/heimdall-home", localPathHeimdall, []uint{}, false, "", false, "versionchecker", true)
		if err != nil {
			utils.WriteError("Error running image:" + err.Error())
			return false
		} else {
			version, err := utils.ExtractVersion(output)
			if err != nil {
				utils.WriteError("Error executing heimdallcli:" + err.Error())
				return false
			}
			fmt.Print(version)
		}

		appstate.UpdateState(appstate.ConfiguringHeimdall)
	}

	if appstate.CurrentState <= appstate.ConfiguringHeimdall {
		// init heimdall
		err := os.RemoveAll(localPathHeimdall) // clear config if any
		if err != nil {
			utils.WriteError("Error during heimdall config:" + err.Error())
			return false
		}
		err = os.MkdirAll(localPathHeimdall, os.ModePerm)
		if err != nil {
			utils.WriteError("Error during heimdall config:" + err.Error())
			return false
		}
		output, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"init", "--home=/heimdall-home", "--chain=mainnet"}, "/heimdall-home", localPathHeimdall, []uint{}, false, "", false, "initializer", true)
		if err != nil {
			utils.WriteError("Error during heimdall init:" + err.Error())
			return false
		} else {
			fmt.Print(output)
		}

		// configure
		// for "laddr" we only need to replace the one in rpc section which is the first one
		err = utils.ReplaceValuesInFile(path.Join(localPathHeimdall, "config", "config.toml"), map[string]string{"moniker": "keepix-node", "laddr": "tcp://0.0.0.0:26657"})
		if err != nil {
			utils.WriteError("Error during heimdall configure:" + err.Error())
			return false
		}

		err = utils.ReplaceValuesInFile(path.Join(localPathHeimdall, "config", "heimdall-config.toml"), map[string]string{"eth_rpc_url": ethereumRPC, "bor_rpc_url": "http://bor:8545"})
		if err != nil {
			utils.WriteError("Error during heimdall configure:" + err.Error())
			return false
		}

		appstate.UpdateState(appstate.ConfiguringBor)
	}

	if appstate.CurrentState <= appstate.ConfiguringBor {
		err := os.RemoveAll(localPathBor) // clear config if any
		if err != nil {
			utils.WriteError("Error during bor config:" + err.Error())
			return false
		}
		err = os.MkdirAll(localPathBor, os.ModePerm)
		if err != nil {
			utils.WriteError("Error during bor config:" + err.Error())
			return false
		}

		// write genesis to file
		genesisFile := path.Join(localPathBor, "genesis.json")
		err = os.WriteFile(genesisFile, []byte(genesis), fs.FileMode(0644))
		if err != nil {
			utils.WriteError("Error writing genesis file: " + err.Error())
			return false
		}

		// write config to file
		tomlFile := path.Join(localPathBor, "config.toml")
		err = os.WriteFile(tomlFile, []byte(configToml), fs.FileMode(0644))
		if err != nil {
			utils.WriteError("Error writing toml file: " + err.Error())
			return false
		}

		appstate.UpdateState(appstate.ConfiguringNetwork)
	}

	if appstate.CurrentState <= appstate.ConfiguringNetwork {
		// recreate the network
		utils.RemoveDockerNetwork("polygon")
		// create docker network
		err := utils.CreateDockerNetwork("polygon")
		if err != nil {
			utils.WriteError("Error creating docker network:" + err.Error())
			return false
		}

		appstate.UpdateState(appstate.NodeInstalled)
	}

	return true
}

func startTask(args []string) bool {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return false
	}
	storage, _ := appstate.GetStoragePath()
	localPathHeimdall := path.Join(storage, "data", "heimdall")
	localPathBor := path.Join(storage, "data", "bor")

	if appstate.CurrentState <= appstate.StartingHeimdall {
		appstate.UpdateState(appstate.StartingHeimdall)
		output, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"start", "--home=/heimdall-home"}, "/heimdall-home", localPathHeimdall, []uint{26657, 26656}, true, "polygon", true, "heimdall", false)
		if err != nil {
			utils.WriteError("Error during heimdall start:" + err.Error())
			return false
		} else {
			fmt.Print(output)
			appstate.UpdateState(appstate.StartingRestServer)
		}
	}

	if appstate.CurrentState <= appstate.StartingRestServer {
		output, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"rest-server", "--home=/heimdall-home", "--node=tcp://heimdall:26657"}, "/heimdall-home", localPathHeimdall, []uint{1317}, true, "polygon", true, "heimdall-rest", false)
		if err != nil {
			utils.WriteError("Error during heimdall rest server start:" + err.Error())
			return false
		} else {
			fmt.Print(output)
			appstate.UpdateState(appstate.StartingBor)
		}
	}

	if appstate.CurrentState <= appstate.StartingBor {
		output, err := utils.DockerRun("0xpolygon/bor:1.1.0", []string{"server", "--datadir=/bor-home", "--config=/bor-home/config.toml"}, "/bor-home", localPathBor, []uint{30303, 8545}, true, "polygon", true, "bor", false)
		if err != nil {
			utils.WriteError("Error during heimdall rest server start:" + err.Error())
			return false
		} else {
			fmt.Print(output)
			appstate.UpdateState(appstate.StartingBor)
		}
	}

	appstate.UpdateState(appstate.NodeStarted)

	return true
}

func stopTask(args []string) bool {
	appstate.UpdateState(appstate.NodeInstalled)
	return true
}

// uninstallTask is an example task for uninstallation purposes
func uninstallTask(args []string) bool {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return false
	}

	err := utils.RemoveImage("0xpolygon/heimdall:1.0.3")
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error response from daemon: No such image") {
			// image already removed
			fmt.Print(err.Error())
		} else {
			utils.WriteError("Error removing image:" + err.Error())
			return false
		}
	}

	err = utils.RemoveImage("0xpolygon/bor:1.1.0")
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error response from daemon: No such image") {
			// image already removed
			fmt.Print(err.Error())
		} else {
			utils.WriteError("Error removing image:" + err.Error())
			return false
		}
	}

	err = utils.RemoveDockerNetwork("polygon")
	if err != nil {
		utils.WriteError("Error removing docker network:" + err.Error())
		return false
	}
	err = appstate.UpdateState(appstate.NoState)
	if err != nil {
		utils.WriteError("Error updating state:" + err.Error())
		return false
	}

	return true
}

// returns plugins installation status
func installedTask(args []string) bool {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return false
	}
	fmt.Print(string(appstate.CurrentState))
	return true
}

type NodeStatus struct {
	NodeState    string `json:"NodeState"`
	Alive        bool   `json:"Alive"`
	IsRegistered bool   `json:"IsRegistered"`
}

// returns plugins status
func statusTask(args []string) bool {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return false
	}
	// Create an instance of NodeStatus
	status := NodeStatus{
		NodeState:    string(appstate.CurrentState),
		Alive:        false,
		IsRegistered: false,
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return false
	}

	// Print the JSON string
	fmt.Print(string(jsonBytes))
	return true
}

type SyncState struct {
	IsSynced     bool `json:"IsSynced"`
	SyncProgress int  `json:"SyncProgress"`
}

func syncStateTask(args []string) bool {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return false
	}
	state, err := utils.GetNodeStatus()
	if err != nil {
		utils.WriteError("Error getting node status:" + err.Error())
		return false
	}

	blockHeight, _ := strconv.Atoi(state.Result.SyncInfo.LatestBlockHeight)

	status := &SyncState{
		IsSynced:     !state.Result.SyncInfo.CatchingUp,
		SyncProgress: blockHeight,
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return false
	}

	fmt.Print(string(jsonBytes))

	return true
}

// validateRequirements checks if all requirements for a task are met
func ValidateRequirements(taskName string) (bool, []string) {
	var missingRequirements []string
	requirements, exists := TaskRequirements[taskName]
	if !exists {
		// No specific requirements
		return true, missingRequirements
	}

	for _, requirement := range requirements {
		checkFunc, exists := SystemRequirements[requirement]
		if !exists || !checkFunc() {
			// Requirement not met
			missingRequirements = append(missingRequirements, requirement)
		}
	}
	return len(missingRequirements) == 0, missingRequirements
}
