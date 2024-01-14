package tasks

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

// returns plugins installation status
func installedTask(args map[string]string) string {
	fmt.Print(string(appstate.CurrentStateString()))
	return RESULT_SUCCESS
}

type NodeStatus struct {
	NodeState string `json:"NodeState"`
	Alive     bool   `json:"Alive"`
}

// returns plugins status
func statusTask(args map[string]string) string {
	_, err := utils.GetHeimdallNodeStatus()
	_, err2 := utils.GetErigonSyncingStatus()

	if !appstate.CurrentState.HeimdallSnapshotDownloaded {
		_, err = utils.SnapshotProgress()
	}

	// Create an instance of NodeStatus
	status := NodeStatus{
		NodeState: appstate.CurrentStateString(),
		Alive:     err == nil && err2 == nil,
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
}

type LogsResponse struct {
	HeimdallLogs string `json:"heimdallLogs"`
	ErigonLogs   string `json:"erigonLogs"`
}

func logsTask(args map[string]string) string {
	erigonLogs := args["erigon"]
	heimdallLogs := args["heimdall"]
	linesAmount, err := strconv.Atoi(args["lines"])
	if err != nil {
		utils.WriteError("Invalid lines amount")
		return RESULT_ERROR
	}

	var logsResponse LogsResponse = LogsResponse{}
	if erigonLogs != "true" && heimdallLogs != "true" {
		return RESULT_SUCCESS
	}
	if erigonLogs == "true" {
		output, err := utils.FetchContainerLogs("erigon", linesAmount)
		if err != nil {
			utils.WriteError("Error getting logs:" + err.Error())
			return RESULT_ERROR
		}
		logsResponse.ErigonLogs = output
	}
	if heimdallLogs == "true" {
		output, err := utils.FetchContainerLogs("heimdall", linesAmount)
		if err != nil {
			utils.WriteError("Error getting logs:" + err.Error())
			return RESULT_ERROR
		}
		logsResponse.HeimdallLogs = output
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(logsResponse)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
}

type SyncState struct {
	IsSynced                bool    `json:"IsSynced"`
	ErigonSyncProgress      float32 `json:"erigonSyncProgress"`
	HeimdallSyncProgress    float32 `json:"heimdallSyncProgress"`
	ErigonStepDescription   string  `json:"erigonStepDescription"`
	HeimdallStepDescription string  `json:"heimdallStepDescription"`
}

func getChainTask(args map[string]string) string {
	if !appstate.CurrentState.IsTestnet {
		return "mainnet"
	} else {
		return "testnet"
	}
}

func syncStateTask(args map[string]string) string {

	erigonState, err := utils.GetErigonSyncingStatus()
	if err != nil {
		utils.WriteError("Error getting erigon node status:" + err.Error())
		return RESULT_ERROR
	}

	var heimdallStepDescription string
	var progress float32
	var heimdallSynced = false

	if appstate.CurrentState.HeimdallSnapshotDownloaded {
		heimdallState, err := utils.GetHeimdallNodeStatus()
		if err != nil {
			utils.WriteError("Error getting heimdall node status:" + err.Error())
			return RESULT_ERROR
		}

		blockHeight, _ := strconv.Atoi(heimdallState.Result.SyncInfo.LatestBlockHeight)
		currentBlockHeight, _ := strconv.Atoi(heimdallState.Result.SyncInfo.CurrentBlockHeight)

		if !heimdallState.Result.SyncInfo.CatchingUp {
			progress = 100
			heimdallStepDescription = "Synced"
		} else {
			progress = float32(blockHeight) / float32(currentBlockHeight) * 100
			heimdallStepDescription = "Syncing"
		}

		if blockHeight == 0 {
			heimdallStepDescription = "Waiting for peers"
		}

		heimdallSynced = !heimdallState.Result.SyncInfo.CatchingUp
	} else {
		heimdallStepDescription = "Downloading snapshot"

		progress, err = utils.SnapshotProgress()
		if err != nil {
			utils.WriteError("Error getting snapshot progress:" + err.Error())
			return RESULT_ERROR
		}
		if progress == 100 {
			heimdallStepDescription = "Snapshot downloaded, restart Heimdall"
		}
	}

	status := &SyncState{
		IsSynced:                heimdallSynced && erigonState.Stage == "Synced",
		ErigonSyncProgress:      erigonState.Progress,
		HeimdallSyncProgress:    progress,
		ErigonStepDescription:   erigonState.Stage,
		HeimdallStepDescription: heimdallStepDescription,
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
}
