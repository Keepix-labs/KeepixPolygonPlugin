package tasks

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/utils"
	"path"
)

func startTask(args map[string]string) string {
	storage, _ := appstate.GetStoragePath()
	localPathHeimdall := path.Join(storage, "data", "heimdall")
	localPathErigon := path.Join(storage, "data", "erigon")

	if appstate.CurrentState.State <= appstate.StartingHeimdall {
		appstate.UpdateState(appstate.StartingHeimdall)
		_ = utils.StopContainerByName("heimdall") // try and stop heimdall if it's already running
		_, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"start", "--home=/heimdall-home"}, "/heimdall-home", localPathHeimdall, []uint{26657, 26656}, true, "polygon", true, "heimdall", false)
		if err != nil {
			utils.WriteError("Error during heimdall start:" + err.Error())
			return RESULT_ERROR
		} else {
			//fmt.Print(output)
			appstate.UpdateState(appstate.StartingRestServer)
		}
	}

	if appstate.CurrentState.State <= appstate.StartingRestServer {
		_ = utils.StopContainerByName("heimdall-rest") // try and stop heimdall-rest if it's already running
		_, err := utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"rest-server", "--home=/heimdall-home", "--node=tcp://heimdall:26657"}, "/heimdall-home", localPathHeimdall, []uint{1317}, true, "polygon", true, "heimdall-rest", false)
		if err != nil {
			utils.WriteError("Error during heimdall rest server start:" + err.Error())
			return RESULT_ERROR
		} else {
			//fmt.Print(output)
			appstate.UpdateState(appstate.StartingErigon)
		}
	}

	if appstate.CurrentState.State <= appstate.StartingErigon {
		chainArg := "--chain=bor-mainnet"
		if appstate.CurrentState.IsTestnet {
			chainArg = "--chain=mumbai"
		}
		_ = utils.StopContainerByName("erigon") // try and stop erigon if it's already running
		_, err := utils.DockerRun("thorax/erigon:v2.53.4", []string{"--datadir=/erigon-home", "--bor.heimdall=http://heimdall-rest:1317", "--private.api.addr=0.0.0.0:9090", "--http.addr=0.0.0.0", chainArg}, "/erigon-home", localPathErigon, []uint{30303, 8545, 9090}, true, "polygon", true, "erigon", false)
		if err != nil {
			utils.WriteError("Error during erigon start:" + err.Error())
			return RESULT_ERROR
		} else {
			//fmt.Print(output)
		}
	}

	appstate.UpdateState(appstate.NodeStarted)

	return RESULT_SUCCESS
}

func stopTask(args map[string]string) string {
	err1 := utils.StopContainerByName("heimdall")
	err2 := utils.StopContainerByName("heimdall-rest")
	err3 := utils.StopContainerByName("erigon")

	if err1 != nil {
		utils.WriteError("Error stopping heimdall:" + err1.Error())
	}
	if err2 != nil {
		utils.WriteError("Error stopping heimdall-rest:" + err2.Error())
	}
	if err3 != nil {
		utils.WriteError("Error stopping erigon:" + err3.Error())
	}

	if err1 != nil || err2 != nil || err3 != nil {
		return RESULT_ERROR
	}
	appstate.UpdateState(appstate.NodeInstalled)
	return RESULT_SUCCESS
}

func resyncTask(args map[string]string) string {
	resyncErigon := args["erigon"]
	resyncHeimdall := args["heimdall"]
	if resyncErigon != "true" && resyncHeimdall != "true" {
		return RESULT_SUCCESS
	}
	res := stopTask(map[string]string{})
	if res != RESULT_SUCCESS {
		utils.WriteError("Error stopping node")
		return RESULT_ERROR
	}
	resRemove := removeData(resyncErigon == "true", resyncHeimdall == "true", false)
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
