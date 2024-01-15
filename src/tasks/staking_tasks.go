package tasks

import (
	"KeepixPlugin/appstate"
	"KeepixPlugin/utils"
	"encoding/json"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type Validator struct {
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	LogoUrl           string  `json:"logoUrl"`
	CommissionPercent int     `json:"commissionPercent"`
	TotalStaked       float64 `json:"totalStaked"`
	Status            string  `json:"status"`
	ContractAddress   string  `json:"contractAddress"`
	UptimePercent     float32 `json:"uptimePercent"`
	DelegationEnabled bool    `json:"delegationEnabled"`
	PerformanceIndex  float32 `json:"performanceIndex"`
	CurrentState      string  `json:"currentState"`
	UserStake         string  `json:"userStake"`
	UserReward        string  `json:"userReward"`
	MinStake          string  `json:"minStake"`
}

type ValidatorsResponse struct {
	Summary struct {
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
		SortBy    string `json:"sortBy"`
		Direction string `json:"direction"`
		Total     int    `json:"total"`
		Size      int    `json:"size"`
	} `json:"summary"`
	Success bool        `json:"success"`
	Status  string      `json:"status"`
	Result  []Validator `json:"result"`
}

const tokenABI = `[
    {
        "constant": false,
        "inputs": [
            {
                "name": "spender",
                "type": "address"
            },
            {
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "approve",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }
]`

const validatorABI = `[
    {
        "constant": false,
        "inputs": [
            {
                "name": "_amount",
                "type": "uint256"
            },
            {
                "name": "_minSharesToMint",
                "type": "uint256"
            }
        ],
        "name": "buyVoucher",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "constant": false,
        "inputs": [
            {
                "name": "claimAmount",
                "type": "uint256"
            },
            {
                "name": "maximumSharesToBurn",
                "type": "uint256"
            }
        ],
        "name": "sellVoucher_new",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    },
	{
        "constant": true,
        "inputs": [
            {
                "name": "user",
                "type": "address"
            }
        ],
        "name": "getTotalStake",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            },
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
	{
        "constant": true,
        "inputs": [],
        "name": "minAmount",
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
	{
        "constant": true,
        "inputs": [],
        "name": "getLiquidRewards",
		"inputs": [
            {
                "name": "user",
                "type": "address"
            }
        ],
        "outputs": [
            {
                "name": "",
                "type": "uint256"
            }
        ],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    },
	{
        "constant": false,
        "name": "withdrawRewards",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }
]`

// weiToEther converts a wei amount (as big.Int) to ether (as float64) with higher precision.
func weiToEther(wei *big.Int) string {
	// Use a big.Float with higher precision
	const precision = 256
	ether := new(big.Float).SetPrec(precision).SetInt(wei)
	divisor := new(big.Float).SetPrec(precision).SetFloat64(1e18)

	// Perform the division
	ether.Quo(ether, divisor)

	// Convert the result to float64
	return ether.Text('f', 18)
}

// fetchValidators fetches validators data from the provided URL and unmarshals into ValidatorsResponse struct.
func fetchValidators() (*ValidatorsResponse, error) {
	var url string
	if appstate.CurrentState.IsTestnet {
		url = "https://staking-api-testnet.polygon.technology/api/v2/validators?limit=10&offset=0&sortBy=delegatedStake"
	} else {
		url = "https://staking-api.polygon.technology/api/v2/validators?limit=10&offset=0&sortBy=delegatedStake"
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return nil, err
	}

	var validatorsResponse ValidatorsResponse
	err = json.Unmarshal(body, &validatorsResponse)
	if err != nil {
		return nil, err
	}

	// add information about min stake and user stake
	client, err := utils.NewBlockchainClient(appstate.CurrentState.RPC)
	if err != nil {
		utils.WriteError("Error creating blockchain client:" + err.Error())
		return nil, err
	}

	addr := common.HexToAddress(appstate.CurrentState.Wallet.Address)

	for index, validator := range validatorsResponse.Result {
		minAmountResult, err := client.CallReadOnlyFunction(validator.ContractAddress, validatorABI, "minAmount")
		if err != nil {
			utils.WriteError("Error calling minAmount:" + err.Error())
			return nil, err
		}
		userStakeResult, err := client.CallReadOnlyFunction(validator.ContractAddress, validatorABI, "getTotalStake", addr)
		if err != nil {
			utils.WriteError("Error calling getTotalStake:" + err.Error())
			return nil, err
		}
		userRewardResult, err := client.CallReadOnlyFunction(validator.ContractAddress, validatorABI, "getLiquidRewards", addr)
		if err != nil {
			utils.WriteError("Error calling getLiquidRewards:" + err.Error())
			return nil, err
		}
		minAmount, ok := minAmountResult[0].(*big.Int)
		if !ok {
			utils.WriteError("Error converting result to bytes")
			return nil, err
		}

		userStake, ok := userStakeResult[0].(*big.Int)
		if !ok {
			utils.WriteError("Error converting result to bytes")
			return nil, err
		}

		userReward, ok := userRewardResult[0].(*big.Int)
		if !ok {
			utils.WriteError("Error converting result to bytes")
			return nil, err
		}

		validatorsResponse.Result[index].MinStake = weiToEther(minAmount)
		validatorsResponse.Result[index].UserStake = weiToEther(userStake)
		validatorsResponse.Result[index].UserReward = weiToEther(userReward)
	}

	return &validatorsResponse, nil
}

// Struct to match the innermost objects in the "result" array.
type Delegate struct {
	BondedValidator   int     `json:"bondedValidator"`
	Stake             float64 `json:"stake"` // Assuming stake is a float64, adjust if necessary
	Address           string  `json:"address"`
	ClaimedReward     float64 `json:"claimedReward"` // Assuming it fits into an int64
	Shares            string  `json:"shares"`
	DeactivationEpoch string  `json:"deactivationEpoch"`
}

// Main struct to match the top level of the JSON.
type DelegatorsResponse struct {
	Success bool       `json:"success"`
	Status  string     `json:"status"`
	Result  []Delegate `json:"result"`
}

// poolsFetchTask fetches the list of validators
func poolsFetchTask(args map[string]string) string {
	response, err := fetchValidators()
	if err != nil {
		utils.WriteError("Error fetching validators:" + err.Error())
		return RESULT_ERROR
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(response.Result)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
}

// unstakeTask unstakes an amount from a validator
func unstakeTask(args map[string]string) string {
	address := args["address"]
	amount := args["amount"]

	if !common.IsHexAddress(address) {
		utils.WriteError("Not a valid hex address")
		return RESULT_ERROR
	}

	client, err := utils.NewBlockchainClient(appstate.CurrentState.RPC)
	if err != nil {
		utils.WriteError("Error creating blockchain client:" + err.Error())
		return RESULT_ERROR
	}

	bigIntAmount := new(big.Int)
	_, success := bigIntAmount.SetString(amount, 10)
	if !success {
		utils.WriteError("Error converting amount to big.Int")
		return RESULT_ERROR
	}

	privateKey, err := utils.ConvertBase64ToPrivateKey(appstate.CurrentState.Wallet.PK)
	if err != nil {
		utils.WriteError("Error converting private key:" + err.Error())
		return RESULT_ERROR
	}

	hash, err := client.ExecuteWriteFunction(privateKey, address, validatorABI, "sellVoucher_new", 300000, bigIntAmount, bigIntAmount)
	if err != nil {
		utils.WriteError("Error executing unstake:" + err.Error())
		return RESULT_ERROR
	}
	receipt, err := client.WaitForTransactionReceipt(hash)
	if err != nil {
		utils.WriteError("Error waiting for receipt:" + err.Error())
		return RESULT_ERROR
	}
	if receipt.Status == 0 {
		utils.WriteError("Transaction failed: " + receipt.TxHash.String())
		return RESULT_ERROR
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(receipt)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
}

// stakeTask stakes an amount on a validator
func stakeTask(args map[string]string) string {
	address := args["address"]
	amount := args["amount"]

	if !common.IsHexAddress(address) {
		utils.WriteError("Not a valid hex address")
		return RESULT_ERROR
	}

	client, err := utils.NewBlockchainClient(appstate.CurrentState.RPC)
	if err != nil {
		utils.WriteError("Error creating blockchain client:" + err.Error())
		return RESULT_ERROR
	}

	bigIntAmount := new(big.Int)
	_, success := bigIntAmount.SetString(amount, 10)
	if !success {
		utils.WriteError("Error converting amount to big.Int")
		return RESULT_ERROR
	}

	zero := big.NewInt(0)

	// execute approval function
	maticAddress := MATIC_ADDR
	if appstate.CurrentState.IsTestnet {
		maticAddress = TESTNET_MATIC_ADDR
	}

	privateKey, err := utils.ConvertBase64ToPrivateKey(appstate.CurrentState.Wallet.PK)
	if err != nil {
		utils.WriteError("Error converting private key:" + err.Error())
		return RESULT_ERROR
	}

	hash, err := client.ExecuteWriteFunction(privateKey, maticAddress, tokenABI, "approve", 80000, common.HexToAddress(address), bigIntAmount)
	if err != nil {
		utils.WriteError("Error executing approval:" + err.Error())
		return RESULT_ERROR
	}
	receipt, err := client.WaitForTransactionReceipt(hash)
	if err != nil {
		utils.WriteError("Error waiting for receipt:" + err.Error())
		return RESULT_ERROR
	}
	if receipt.Status == 0 {
		utils.WriteError("Approval transaction failed: " + receipt.TxHash.String())
		return RESULT_ERROR
	}

	hash, err = client.ExecuteWriteFunction(privateKey, address, validatorABI, "buyVoucher", 300000, bigIntAmount, zero)
	if err != nil {
		utils.WriteError("Error executing stake:" + err.Error())
		return RESULT_ERROR
	}

	receipt, err = client.WaitForTransactionReceipt(hash)
	if err != nil {
		utils.WriteError("Error waiting for receipt:" + err.Error())
		return RESULT_ERROR
	}
	if receipt.Status == 0 {
		utils.WriteError("Staking transaction failed: " + receipt.TxHash.String())
		return RESULT_ERROR
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(receipt)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
}

// rewardTask gets the reward from a validator
func rewardTask(args map[string]string) string {
	address := args["address"]

	if !common.IsHexAddress(address) {
		utils.WriteError("Not a valid hex address")
		return RESULT_ERROR
	}

	client, err := utils.NewBlockchainClient(appstate.CurrentState.RPC)
	if err != nil {
		utils.WriteError("Error creating blockchain client:" + err.Error())
		return RESULT_ERROR
	}

	privateKey, err := utils.ConvertBase64ToPrivateKey(appstate.CurrentState.Wallet.PK)
	if err != nil {
		utils.WriteError("Error converting private key:" + err.Error())
		return RESULT_ERROR
	}

	hash, err := client.ExecuteWriteFunction(privateKey, address, validatorABI, "withdrawRewards", 180000)
	if err != nil {
		utils.WriteError("Error executing rewards:" + err.Error())
		return RESULT_ERROR
	}

	receipt, err := client.WaitForTransactionReceipt(hash)
	if err != nil {
		utils.WriteError("Error waiting for receipt:" + err.Error())
		return RESULT_ERROR
	}
	if receipt.Status == 0 {
		utils.WriteError("Reward claiming transaction failed: " + receipt.TxHash.String())
		return RESULT_ERROR
	}

	// Serialize the struct to JSON
	jsonBytes, err := json.Marshal(receipt)
	if err != nil {
		utils.WriteError("Error serializing to JSON:" + err.Error())
		return RESULT_ERROR
	}

	return string(jsonBytes)
}
