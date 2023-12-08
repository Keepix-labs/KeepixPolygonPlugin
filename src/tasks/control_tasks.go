package tasks

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/utils"
	"fmt"
	"path"
)

func startTask(args map[string]string) string {
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

func stopTask(args map[string]string) string {
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

func resyncTask(args map[string]string) string {
	resyncBor := args["bor"]
	resyncHeimdall := args["heimdall"]
	if resyncBor != "true" && resyncHeimdall != "true" {
		return RESULT_SUCCESS
	}
	res := stopTask(map[string]string{})
	if res != RESULT_SUCCESS {
		utils.WriteError("Error stopping node")
		return RESULT_ERROR
	}
	resRemove := removeData(resyncBor == "true", resyncHeimdall == "true")
	if !resRemove {
		utils.WriteError("Error removing data")
		return RESULT_ERROR
	}
	res = startTask(map[string]string{})
	if res != RESULT_SUCCESS {
		utils.WriteError("Error starting node")
		return RESULT_ERROR
	}
	return RESULT_SUCCESS
}

func restartTask(args map[string]string) string {
	res := stopTask(map[string]string{})
	if res != RESULT_SUCCESS {
		utils.WriteError("Error stopping node")
		return RESULT_ERROR
	}
	res = startTask(map[string]string{})
	if res != RESULT_SUCCESS {
		utils.WriteError("Error starting node")
		return RESULT_ERROR
	}
	return RESULT_SUCCESS
}
