package main

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/tasks"
	"KeepixPlugin/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type AppResult struct {
	Result string `json:"jsonResult"`
	Stdout string `json:"stdOut"`
	Stderr string `json:"stdErr"`
}

var version string

func main() {
	if len(os.Args) >= 2 {
		if os.Args[1] == "--version" {
			fmt.Print(version)
			os.Exit(0)
		}
	}
	appResult, err := App()
	if err != nil {
		fmt.Print("Error running the application:", err)
		os.Exit(1)
	}

	// Convert the result to JSON
	jsonResult, err := json.Marshal(appResult)
	if err != nil {
		fmt.Println("Error marshalling result to JSON:", err)
		return
	}

	fmt.Print(string(jsonResult))
}

// App runs the application and captures stdout and stderr.
func App() (AppResult, error) {
	// Backup original stdout and stderr
	origStdout := os.Stdout
	origStderr := os.Stderr
	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	// Create pipes to capture stdout and stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	// Run the application logic
	result := Plugin() // Replace with your application logic

	// Close the writers and capture the output
	wOut.Close()
	wErr.Close()
	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)

	return AppResult{
		Result: result,
		Stdout: bufOut.String(),
		Stderr: bufErr.String(),
	}, nil
}

func Plugin() string {
	if len(os.Args) != 2 {
		return tasks.RESULT_ERROR
	}

	var input struct {
		Key string `json:"key"`
	}

	err := appstate.LoadState()
	if err != nil {
		utils.WriteError(err.Error())
		return tasks.RESULT_ERROR
	} else {
		err = json.Unmarshal([]byte(os.Args[1]), &input)
		if err != nil {
			utils.WriteError(err.Error())
			return tasks.RESULT_ERROR
		} else {
			taskFunc, exists := tasks.TaskMap[input.Key]
			if !exists {
				utils.WriteError("Invalid command")
				return tasks.RESULT_ERROR
			} else {
				// Parse arguments
				var dataMap map[string]string
				if err := json.Unmarshal([]byte(os.Args[1]), &dataMap); err != nil {
					utils.WriteError("Invalid args")
					return tasks.RESULT_ERROR
				}
				// remove the key arg
				delete(dataMap, "key")
				validated, missing := tasks.ValidateRequirements(input.Key)
				if !validated {
					utils.WriteError("Missing requirements for command: " + strings.Join(missing, ", "))
					return tasks.RESULT_ERROR
				}
				validated, missing = tasks.ValidateArgs(input.Key, dataMap)
				if !validated {
					utils.WriteError("Missing arguments for command: " + strings.Join(missing, ", "))
					return tasks.RESULT_ERROR
				}
				return taskFunc(dataMap)
			}
		}
	}
}
