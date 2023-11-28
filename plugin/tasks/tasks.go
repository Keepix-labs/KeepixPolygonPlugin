package tasks

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/utils"
	"encoding/json"
	"fmt"
)

// TaskFunc defines the signature for task functions
type TaskFunc func([]string) bool

// TaskMap maps task names to their corresponding functions
var TaskMap = map[string]TaskFunc{
	"install":   installTask,
	"uninstall": uninstallTask,
	"installed": installedTask,
	"status":    statusTask,
	// Add other tasks here...
}

// TaskRequirements maps task names to their required system conditions
var TaskRequirements = map[string][]string{
	"install":   {"docker", "uninstalled", "linux", "cpu4"},
	"uninstall": {"docker", "installed"},
	"installed": {"docker"},
	"status":    {},
}

// installTask is an example task for installation purposes
func installTask(args []string) bool {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return false
	}

	err := appstate.UpdateState(appstate.InstallingNode)
	if err != nil {
		utils.WriteError("Error updating state:" + err.Error())
		return false
	}

	// Implement your task logic here using the arguments
	err = utils.PullImage("0xpolygon/heimdall:latest")
	if err != nil {
		utils.WriteError("Error pulling image:" + err.Error())
		return false
	}

	err = appstate.UpdateState(appstate.NodeStopped)
	if err != nil {
		utils.WriteError("Error updating state:" + err.Error())
		return false
	}

	return true
}

// uninstallTask is an example task for uninstallation purposes
func uninstallTask(args []string) bool {
	if len(args) > 0 {
		utils.WriteError("Too many arguments")
		return false
	}

	err := utils.RemoveImage("0xpolygon/heimdall:latest")
	if err != nil {
		utils.WriteError("Error removing image:" + err.Error())
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

// Define a struct to hold your data
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

	// Print or return the JSON string
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
