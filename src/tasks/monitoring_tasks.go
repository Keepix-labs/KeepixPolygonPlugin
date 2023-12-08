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
	NodeState    string `json:"NodeState"`
	Alive        bool   `json:"Alive"`
	IsRegistered bool   `json:"IsRegistered"`
}

// returns plugins status
func statusTask(args map[string]string) string {
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

type LogsResponse struct {
	HeimdallLogs string `json:"heimdallLogs"`
	BorLogs      string `json:"borLogs"`
}

func logsTask(args map[string]string) string {
	bogLogs := args["bor"]
	heimdallLogs := args["heimdall"]
	linesAmount, err := strconv.Atoi(args["lines"])
	if err != nil {
		utils.WriteError("Invalid lines amount")
		return RESULT_ERROR
	}

	var logsResponse LogsResponse = LogsResponse{}
	if bogLogs != "true" && heimdallLogs != "true" {
		return RESULT_SUCCESS
	}
	if bogLogs == "true" {
		output, err := utils.FetchContainerLogs("bor", linesAmount)
		if err != nil {
			utils.WriteError("Error getting logs:" + err.Error())
			return RESULT_ERROR
		}
		logsResponse.BorLogs = output
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
	BorSyncProgress         float32 `json:"borSyncProgress"`
	HeimdallSyncProgress    float32 `json:"heimdallSyncProgress"`
	BorStepDescription      string  `json:"borStepDescription"`
	HeimdallStepDescription string  `json:"heimdallStepDescription"`
}

func syncStateTask(args map[string]string) string {
	heimdallState, err := utils.GetHeimdallNodeStatus()
	if err != nil {
		utils.WriteError("Error getting heimdall node status:" + err.Error())
		return RESULT_ERROR
	}

	borState, err := utils.GetBorNodeStatus()
	if err != nil {
		utils.WriteError("Error getting bor node status:" + err.Error())
		return RESULT_ERROR
	}

	var borStepDescription string
	var heimdallStepDescription string

	var progress float32
	if !borState.CatchingUp {
		progress = 100
		borStepDescription = "Synced"
	} else {
		highestBlock, _ := strconv.ParseInt(borState.Result.HighestBlock, 16, 64)
		currentBlock, _ := strconv.ParseInt(borState.Result.CurrentBlock, 16, 64)

		if highestBlock == 0 {
			progress = 0
			borStepDescription = "Waiting for heimdall sync"
		} else {
			progress = float32(currentBlock) / float32(highestBlock) * 100
			borStepDescription = "Syncing"
		}
	}

	blockHeight, _ := strconv.Atoi(heimdallState.Result.SyncInfo.LatestBlockHeight)
	currentBlockHeight, _ := strconv.Atoi(heimdallState.Result.SyncInfo.CurrentBlockHeight)

	var progress2 float32
	if !heimdallState.Result.SyncInfo.CatchingUp {
		progress2 = 100
		heimdallStepDescription = "Synced"
	} else {
		progress2 = float32(blockHeight) / float32(currentBlockHeight) * 100
		heimdallStepDescription = "Syncing"
		// override bor status for a clearer view of what's happening
		borStepDescription = "Waiting for heimdall sync"
		progress = 0
	}

	if blockHeight == 0 {
		heimdallStepDescription = "Waiting for peers"
	}

	status := &SyncState{
		IsSynced:                !heimdallState.Result.SyncInfo.CatchingUp && !borState.CatchingUp,
		BorSyncProgress:         progress,
		HeimdallSyncProgress:    progress2,
		BorStepDescription:      borStepDescription,
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
