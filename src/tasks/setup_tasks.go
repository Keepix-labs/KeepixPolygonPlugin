package tasks

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/utils"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path"
)

//go:embed conf/heimdall/config.toml
var configHeimdallToml string

// installTask is an example task for installation purposes
func installTask(args map[string]string) string {
	ethereumRPC := args["ethereumRPC"]
	if !utils.IsValidURL(ethereumRPC) {
		utils.WriteError("Invalid ethereumRPC")
		return RESULT_ERROR
	}

	isTestnet := args["testnet"] == "true"
	isAutoStart := args["autostart"] == "true"
	mnemonic := args["mnemonic"]

	appstate.UpdateChain(isTestnet)
	appstate.UpdateRPC(ethereumRPC)
	err := utils.LoadAccountFromMnemonic(mnemonic)
	if err != nil {
		utils.WriteError("Error loading account from mnemonic:" + err.Error())
		return RESULT_ERROR
	}

	storage, _ := appstate.GetStoragePath()
	localPathHeimdall := path.Join(storage, "data", "heimdall")
	localPathErigon := path.Join(storage, "data", "erigon")

	if appstate.CurrentState.State <= appstate.InstallingNode {
		// not installed yet
		fmt.Println("Installing node")
		appstate.UpdateState(appstate.InstallingNode)

		err := utils.PullImage("0xpolygon/heimdall:1.0.3")
		if err != nil {
			utils.WriteError("Error pulling heimdall image:" + err.Error())
			return RESULT_ERROR
		}

		err = utils.PullImage("thorax/erigon:v2.53.4")
		if err != nil {
			utils.WriteError("Error pulling erigon image:" + err.Error())
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

		fmt.Println("Successfully installed heimdall")
		appstate.UpdateState(appstate.ConfiguringHeimdall)
	}

	if appstate.CurrentState.State <= appstate.ConfiguringHeimdall {
		fmt.Println("Configuring heimdall...")
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
		chainArg := "--chain=mainnet"
		if isTestnet {
			chainArg = "--chain=mumbai"
			fmt.Println("Configuring heimdall for mumbai testnet")
		} else {
			fmt.Println("Configuring heimdall for mainnet")
		}
		_, err = utils.DockerRun("0xpolygon/heimdall:1.0.3", []string{"init", "--home=/heimdall-home", chainArg}, "/heimdall-home", localPathHeimdall, []uint{}, false, "", false, "initializer", true)
		if err != nil {
			utils.WriteError("Error during heimdall init:" + err.Error())
			return RESULT_ERROR
		} else {
			fmt.Println("Successfully initialized heimdall conf")
		}

		// configure
		// write config to file
		err = os.WriteFile(path.Join(localPathHeimdall, "config", "config.toml"), []byte(configHeimdallToml), fs.FileMode(0644))
		if err != nil {
			utils.WriteError("Error writing toml file: " + err.Error())
			return RESULT_ERROR
		}

		err = utils.ReplaceValuesInFile(path.Join(localPathHeimdall, "config", "heimdall-config.toml"), map[string]string{"eth_rpc_url": ethereumRPC, "bor_rpc_url": "http://erigon:8545"})
		if err != nil {
			utils.WriteError("Error during heimdall configure:" + err.Error())
			return RESULT_ERROR
		}
		fmt.Println("Successfully configured heimdall")
		appstate.UpdateState(appstate.ConfiguringErigon)
	}

	if appstate.CurrentState.State <= appstate.ConfiguringErigon {
		fmt.Println("Configuring erigon...")
		err := os.RemoveAll(localPathErigon) // clear config if any
		if err != nil {
			utils.WriteError("Error during erigon config:" + err.Error())
			return RESULT_ERROR
		}
		err = os.MkdirAll(localPathErigon, os.ModePerm)
		if err != nil {
			utils.WriteError("Error during erigon config:" + err.Error())
			return RESULT_ERROR
		}
		fmt.Println("Successfully configured erigon")
		appstate.UpdateState(appstate.ConfiguringNetwork)
	}

	if appstate.CurrentState.State <= appstate.ConfiguringNetwork {
		// recreate the network
		utils.RemoveDockerNetworkIfExists("polygon")
		// create docker network
		err := utils.CreateDockerNetwork("polygon")
		if err != nil {
			utils.WriteError("Error creating docker network:" + err.Error())
			return RESULT_ERROR
		}
		fmt.Println("Successfully installed node")
		appstate.UpdateState(appstate.NodeInstalled)
	}

	if isAutoStart {
		return startTask(args)
	}

	return RESULT_SUCCESS
}

// removeData removes chain data from erigon and heimdall, if all is true, it removes all data
func removeData(erigon bool, heimdall bool, all bool) bool {
	if !erigon && !heimdall {
		return true
	}
	storage, _ := appstate.GetStoragePath()
	// remove data folders using docker because of permission issues
	folders := ""
	if all {
		folders += "/plugin/data/heimdall /plugin/data/erigon"
	} else {
		if erigon {
			folders += "/plugin/data/erigon/bor /plugin/data/erigon/chaindata "
		}
		if heimdall {
			folders += "/plugin/data/heimdall/data/*.db"
		}
	}
	err := utils.RemoveHostFolderUsingContainer("/plugin", storage, folders)
	if err != nil {
		utils.WriteError("Error removing data:" + err.Error())
		return false
	}
	return true
}

// uninstallTask is an example task for uninstallation purposes
func uninstallTask(args map[string]string) string {
	if !removeData(true, true, true) {
		return RESULT_ERROR
	}

	fmt.Println("Removing images and their containers...")

	err := utils.RemoveImageIfExists("0xpolygon/heimdall:1.0.3")
	if err != nil {
		utils.WriteError("Error removing image:" + err.Error())
		return RESULT_ERROR
	}

	err = utils.RemoveImageIfExists("thorax/erigon:v2.53.4")
	if err != nil {
		utils.WriteError("Error removing image:" + err.Error())
		return RESULT_ERROR
	}

	fmt.Println("Successfully removed docker images")

	err = utils.RemoveDockerNetworkIfExists("polygon")
	if err != nil {
		utils.WriteError("Error removing docker network:" + err.Error())
		return RESULT_ERROR
	}

	fmt.Println("Removing plugin data")
	storage, _ := appstate.GetStoragePath()
	// remove rest of plugin data
	err = os.RemoveAll(storage)
	if err != nil {
		utils.WriteError("Error removing data folder:" + err.Error())
		return RESULT_ERROR
	}
	fmt.Println("Successfully removed plugin data")
	return RESULT_SUCCESS
}
