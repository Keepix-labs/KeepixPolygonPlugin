package tasks

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/utils"
	"encoding/json"
	"fmt"
)

type WalletResponse struct {
	Wallet       string `json:"Wallet"`
	MATICBalance string `json:"maticBalance"`
	ETHBalance   string `json:"ethBalance"`
}

const MATIC_ADDR = "0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0"
const TESTNET_MATIC_ADDR = "0x499d11E0b6eAC7c0593d8Fb292DCBbF815Fb29Ae"

// walletFetchTask fetches the stored wallet data
func walletFetchTask(args map[string]string) string {
	client, err := utils.NewBlockchainClient(appstate.CurrentState.RPC)
	if err != nil {
		utils.WriteError("Error creating blockchain client:" + err.Error())
		return RESULT_ERROR
	}
	maticAddress := MATIC_ADDR
	if appstate.CurrentState.IsTestnet {
		maticAddress = TESTNET_MATIC_ADDR
	}
	maticBalance, err := client.GetERC20Balance(maticAddress, appstate.CurrentState.Wallet.Address)
	if err != nil {
		utils.WriteError("Error fetching MATIC balance:" + err.Error())
		return RESULT_ERROR
	}

	ethBalance, err := client.GetETHBalance(appstate.CurrentState.Wallet.Address)
	if err != nil {
		utils.WriteError("Error fetching ETH balance:" + err.Error())
		return RESULT_ERROR
	}

	wallet := &WalletResponse{
		Wallet:       appstate.CurrentState.Wallet.Address,
		MATICBalance: maticBalance.String(),
		ETHBalance:   ethBalance.String(),
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(wallet)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
}

// walletLoadTask loads a wallet from private key
func walletLoadTask(args map[string]string) string {
	mnemonic := args["mnemonic"]
	privateKey := args["privateKey"]
	if mnemonic == "" && privateKey == "" {
		utils.WriteError("No mnemonic or private key provided")
		return RESULT_ERROR
	}
	if mnemonic != "" && privateKey != "" {
		utils.WriteError("Provide Mnemonic or Private key, not both")
		return RESULT_ERROR
	}
	if mnemonic != "" {
		fmt.Println("Loading wallet from mnemonic...")
		utils.LoadAccountFromMnemonic(mnemonic)
		return RESULT_SUCCESS
	}
	if privateKey != "" {
		fmt.Println("Loading wallet from private key...")
		utils.LoadAccountFromPrivateKey(privateKey)
		return RESULT_SUCCESS
	}
	return RESULT_ERROR
}

// walletPurgeTask removes the stored wallet data
func walletPurgeTask(args map[string]string) string {
	err := appstate.UpdateAccount("", "")
	if err != nil {
		utils.WriteError("Error purging wallet:" + err.Error())
		return RESULT_ERROR
	}
	return RESULT_SUCCESS
}
