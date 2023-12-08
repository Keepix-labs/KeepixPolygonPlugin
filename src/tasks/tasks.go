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
)

// TaskFunc defines the signature for task functions
type TaskFunc func([]string) string

// TaskMap maps task names to their corresponding functions
var TaskMap = map[string]TaskFunc{
	"install":    installTask,
	"uninstall":  uninstallTask,
	"installed":  installedTask,
	"status":     statusTask,
	"start":      startTask,
	"stop":       stopTask,
	"sync-state": syncStateTask,
	"resync":     resyncTask,
	"restart":    restartTask,
}

// TaskRequirements maps task names to their required system conditions
var TaskRequirements = map[string][]string{
	"install":    {"docker", "uninstalled", "linux", "cpu4"},
	"uninstall":  {"docker", "stopped"},
	"installed":  {"docker"},
	"status":     {},
	"start":      {"docker", "stopped"},
	"stop":       {"docker", "running"},
	"sync-state": {"docker", "running"},
	"resync":     {"docker", "installed"},
	"restart":    {"docker", "running"},
}

//go:embed conf/bor/genesis.json
var genesis string

//go:embed conf/bor/config.toml
var configToml string

//go:embed conf/heimdall/config.toml
var configHeimdallToml string

// installTask is an example task for installation purposes
func installTask(args []string) string {
	if len(args) != 1 {
		utils.WriteError("Invalid args")
		return RESULT_ERROR
	}

	ethereumRPC := args[0]
	if !utils.IsValidURL(ethereumRPC) {
		utils.WriteError("Invalid ethereumRPC")
		return RESULT_ERROR
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
			return RESULT_ERROR
		}

		err = utils.PullImage("0xpolygon/bor:1.1.0")
		if err != nil {
			utils.WriteError("Error pulling bor image:" + err.Error())
			return RESULT_ERROR
		}

		// setting up local config path
		err = os.MkdirAll(localPathHeimdall, os.ModePerm)
		if err != nil {
			utils.WriteError("Error creating local path:" + err.Error())
			return RESULT_ERROR
		}

		// check heimdall
		output, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"heimdallcli", "version"}, "/heimdall-home", localPathHeimdall, []uint{}, false, "", false, "versionchecker", true)
		if err != nil {
			utils.WriteError("Error running image:" + err.Error())
			return RESULT_ERROR
		} else {
			version, err := utils.ExtractVersion(output)
			if err != nil {
				utils.WriteError("Error executing heimdallcli:" + err.Error())
				return RESULT_ERROR
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
			return RESULT_ERROR
		}
		err = os.MkdirAll(localPathHeimdall, os.ModePerm)
		if err != nil {
			utils.WriteError("Error during heimdall config:" + err.Error())
			return RESULT_ERROR
		}
		output, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"init", "--home=/heimdall-home", "--chain=mainnet"}, "/heimdall-home", localPathHeimdall, []uint{}, false, "", false, "initializer", true)
		if err != nil {
			utils.WriteError("Error during heimdall init:" + err.Error())
			return RESULT_ERROR
		} else {
			fmt.Print(output)
		}

		// configure
		// write config to file
		err = os.WriteFile(path.Join(localPathHeimdall, "config", "config.toml"), []byte(configHeimdallToml), fs.FileMode(0644))
		if err != nil {
			utils.WriteError("Error writing toml file: " + err.Error())
			return RESULT_ERROR
		}

		err = utils.ReplaceValuesInFile(path.Join(localPathHeimdall, "config", "heimdall-config.toml"), map[string]string{"eth_rpc_url": ethereumRPC, "bor_rpc_url": "http://bor:8545"})
		if err != nil {
			utils.WriteError("Error during heimdall configure:" + err.Error())
			return RESULT_ERROR
		}

		appstate.UpdateState(appstate.ConfiguringBor)
	}

	if appstate.CurrentState <= appstate.ConfiguringBor {
		err := os.RemoveAll(localPathBor) // clear config if any
		if err != nil {
			utils.WriteError("Error during bor config:" + err.Error())
			return RESULT_ERROR
		}
		err = os.MkdirAll(localPathBor, os.ModePerm)
		if err != nil {
			utils.WriteError("Error during bor config:" + err.Error())
			return RESULT_ERROR
		}

		// write genesis to file
		genesisFile := path.Join(localPathBor, "genesis.json")
		err = os.WriteFile(genesisFile, []byte(genesis), fs.FileMode(0644))
		if err != nil {
			utils.WriteError("Error writing genesis file: " + err.Error())
			return RESULT_ERROR
		}

		// write config to file
		tomlFile := path.Join(localPathBor, "config.toml")
		err = os.WriteFile(tomlFile, []byte(configToml), fs.FileMode(0644))
		if err != nil {
			utils.WriteError("Error writing toml file: " + err.Error())
			return RESULT_ERROR
		}

		appstate.UpdateState(appstate.ConfiguringNetwork)
	}

	if appstate.CurrentState <= appstate.ConfiguringNetwork {
		// recreate the network
		utils.RemoveDockerNetworkIfExists("polygon")
		// create docker network
		err := utils.CreateDockerNetwork("polygon")
		if err != nil {
			utils.WriteError("Error creating docker network:" + err.Error())
			return RESULT_ERROR
		}

		appstate.UpdateState(appstate.NodeInstalled)
	}

	return RESULT_SUCCESS
}

func startTask(args []string) string {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return RESULT_ERROR
	}
	storage, _ := appstate.GetStoragePath()
	localPathHeimdall := path.Join(storage, "data", "heimdall")
	localPathBor := path.Join(storage, "data", "bor")

	if appstate.CurrentState <= appstate.StartingHeimdall {
		appstate.UpdateState(appstate.StartingHeimdall)
		output, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"start", "--home=/heimdall-home"}, "/heimdall-home", localPathHeimdall, []uint{26657, 26656}, true, "polygon", true, "heimdall", false)
		if err != nil {
			utils.WriteError("Error during heimdall start:" + err.Error())
			return RESULT_ERROR
		} else {
			fmt.Print(output)
			appstate.UpdateState(appstate.StartingRestServer)
		}
	}

	if appstate.CurrentState <= appstate.StartingRestServer {
		output, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"rest-server", "--home=/heimdall-home", "--node=tcp://heimdall:26657"}, "/heimdall-home", localPathHeimdall, []uint{1317}, true, "polygon", true, "heimdall-rest", false)
		if err != nil {
			utils.WriteError("Error during heimdall rest server start:" + err.Error())
			return RESULT_ERROR
		} else {
			fmt.Print(output)
			appstate.UpdateState(appstate.StartingBor)
		}
	}

	if appstate.CurrentState <= appstate.StartingBor {
		output, err := utils.DockerRun("0xpolygon/bor:1.1.0", []string{"server", "--datadir=/bor-home", "--config=/bor-home/config.toml"}, "/bor-home", localPathBor, []uint{30303, 8545}, true, "polygon", true, "bor", false)
		if err != nil {
			utils.WriteError("Error during heimdall rest server start:" + err.Error())
			return RESULT_ERROR
		} else {
			fmt.Print(output)
			appstate.UpdateState(appstate.StartingBor)
		}
	}

	appstate.UpdateState(appstate.NodeStarted)

	return RESULT_SUCCESS
}

func stopTask(args []string) string {
	err1 := utils.StopContainerByName("heimdall")
	err2 := utils.StopContainerByName("heimdall-rest")
	err3 := utils.StopContainerByName("bor")

	if err1 != nil {
		utils.WriteError("Error stopping heimdall:" + err1.Error())
	}
	if err2 != nil {
		utils.WriteError("Error stopping heimdall-rest:" + err2.Error())
	}
	if err3 != nil {
		utils.WriteError("Error stopping bor:" + err3.Error())
	}

	if err1 != nil || err2 != nil || err3 != nil {
		return RESULT_ERROR
	}
	appstate.UpdateState(appstate.NodeInstalled)
	return RESULT_SUCCESS
}

func removeData(bor bool, heimdall bool) bool {
	if !bor && !heimdall {
		return true
	}
	storage, _ := appstate.GetStoragePath()
	// remove data folders using docker because of permission issues
	folders := ""
	if bor {
		folders += "/plugin/data/bor/bor /plugin/data/bor/keystore "
	}
	if heimdall {
		folders += "/plugin/data/heimdall/data/*.db"
	}
	err := utils.RemoveHostFolderUsingContainer("/plugin", storage, folders)
	if err != nil {
		utils.WriteError("Error removing data:" + err.Error())
		return false
	}
	return true
}

// uninstallTask is an example task for uninstallation purposes
func uninstallTask(args []string) string {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return RESULT_ERROR
	}

	if !removeData(true, true) {
		return RESULT_ERROR
	}

	err := utils.RemoveImageIfExists("0xpolygon/heimdall:1.0.3")
	if err != nil {
		utils.WriteError("Error removing image:" + err.Error())
		return RESULT_ERROR
	}

	err = utils.RemoveImageIfExists("0xpolygon/bor:1.1.0")
	if err != nil {
		utils.WriteError("Error removing image:" + err.Error())
		return RESULT_ERROR
	}

	err = utils.RemoveDockerNetworkIfExists("polygon")
	if err != nil {
		utils.WriteError("Error removing docker network:" + err.Error())
		return RESULT_ERROR
	}
	storage, _ := appstate.GetStoragePath()
	// remove rest of plugin data
	err = os.RemoveAll(storage)
	if err != nil {
		utils.WriteError("Error removing data folder:" + err.Error())
		return RESULT_ERROR
	}
	return RESULT_SUCCESS
}

func resyncTask(args []string) string {
	if len(args) != 2 {
		utils.WriteError("Invalid args")
		return RESULT_ERROR
	}

	resyncBor := args[0]
	resyncHeimdall := args[1]
	if resyncBor != "true" && resyncHeimdall != "true" {
		return RESULT_SUCCESS
	}
	res := stopTask([]string{})
	if res != RESULT_SUCCESS {
		utils.WriteError("Error stopping node")
		return RESULT_ERROR
	}
	resRemove := removeData(resyncBor == "true", resyncHeimdall == "true")
	if !resRemove {
		utils.WriteError("Error removing data")
		return RESULT_ERROR
	}
	res = startTask([]string{})
	if res != RESULT_SUCCESS {
		utils.WriteError("Error starting node")
		return RESULT_ERROR
	}
	return RESULT_SUCCESS
}

func restartTask(args []string) string {
	if len(args) != 0 {
		utils.WriteError("Invalid args")
		return RESULT_ERROR
	}
	res := stopTask([]string{})
	if res != RESULT_SUCCESS {
		utils.WriteError("Error stopping node")
		return RESULT_ERROR
	}
	res = startTask([]string{})
	if res != RESULT_SUCCESS {
		utils.WriteError("Error starting node")
		return RESULT_ERROR
	}
	return RESULT_SUCCESS
}

// returns plugins installation status
func installedTask(args []string) string {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return RESULT_ERROR
	}
	fmt.Print(string(appstate.CurrentStateString()))
	return RESULT_SUCCESS
}

type NodeStatus struct {
	NodeState    string `json:"NodeState"`
	Alive        bool   `json:"Alive"`
	IsRegistered bool   `json:"IsRegistered"`
}

// returns plugins status
func statusTask(args []string) string {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return RESULT_ERROR
	}

	_, err := utils.GetHeimdallNodeStatus()
	_, err2 := utils.GetBorNodeStatus()

	// Create an instance of NodeStatus
	status := NodeStatus{
		NodeState:    appstate.CurrentStateString(),
		Alive:        err == nil && err2 == nil,
		IsRegistered: false,
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
}

type SyncState struct {
	IsSynced             bool    `json:"IsSynced"`
	BorSyncProgress      float32 `json:"borSyncProgress"`
	HeimdallSyncProgress float32 `json:"heimdallSyncProgress"`
}

func syncStateTask(args []string) string {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return RESULT_ERROR
	}
	heimdallState, err := utils.GetHeimdallNodeStatus()
	if err != nil {
		utils.WriteError("Error getting heimdall node status:" + err.Error())
		return RESULT_ERROR
	}

	blockHeight, _ := strconv.Atoi(heimdallState.Result.SyncInfo.LatestBlockHeight)

	borState, err := utils.GetBorNodeStatus()
	if err != nil {
		utils.WriteError("Error getting bor node status:" + err.Error())
		return RESULT_ERROR
	}

	var progress float32
	if !borState.CatchingUp {
		progress = 100
	} else {
		highestBlock, _ := strconv.ParseInt(borState.Result.HighestBlock, 16, 64)
		currentBlock, _ := strconv.ParseInt(borState.Result.CurrentBlock, 16, 64)

		if highestBlock == 0 {
			progress = 0
		} else {
			progress = float32(currentBlock) / float32(highestBlock) * 100
		}
	}

	var progress2 float32
	if !heimdallState.Result.SyncInfo.CatchingUp {
		progress2 = 100
	} else {
		progress2 = float32(blockHeight) //nothing better at the moment
	}

	status := &SyncState{
		IsSynced:             !heimdallState.Result.SyncInfo.CatchingUp && !borState.CatchingUp,
		BorSyncProgress:      progress,
		HeimdallSyncProgress: progress2,
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
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
