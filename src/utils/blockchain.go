package utils

import (
	"KeepixPlugin/appstate"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

// LoadAccountFromMnemonic loads an account from a mnemonic.
func LoadAccountFromMnemonic(mnemonic string) error {
	// Generate a binary seed from the mnemonic.
	seed := bip39.NewSeed(mnemonic, "") // Second parameter is the optional passphrase.

	// Generate a new master key from the seed.
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return err
	}

	// Derive the keys step by step according to the standard Ethereum derivation path m/44'/60'/0'/0/0
	const (
		purpose      = bip32.FirstHardenedChild + 44
		coinType     = bip32.FirstHardenedChild + 60
		account      = bip32.FirstHardenedChild
		change       = 0
		addressIndex = 0
	)

	// Deriving the purpose key
	purposeKey, err := masterKey.NewChildKey(purpose)
	if err != nil {
		return err
	}

	// Deriving the coin type key
	coinTypeKey, err := purposeKey.NewChildKey(coinType)
	if err != nil {
		return err
	}

	// Deriving the account key
	accountKey, err := coinTypeKey.NewChildKey(account)
	if err != nil {
		return err
	}

	// Deriving the change key
	changeKey, err := accountKey.NewChildKey(change)
	if err != nil {
		return err
	}

	// Deriving the address index key
	childKey, err := changeKey.NewChildKey(addressIndex)
	if err != nil {
		return err
	}

	// Get the ECDSA private key from the derived key
	privateKeyECDSA, err := crypto.ToECDSA(childKey.Key)
	if err != nil {
		return err
	}
	// Format the ECDSA private key as a hexadecimal string
	privateKeyHex := hex.EncodeToString(privateKeyECDSA.D.Bytes())

	// Obtain the public key from the private key
	publicKeyECDSA := privateKeyECDSA.Public()
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA.(*ecdsa.PublicKey))

	// Take the Keccak-256 hash of the public key
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:]) // omit the 0x04 prefix

	// The Ethereum address is the last 20 bytes of the hashed public key
	address := crypto.PubkeyToAddress(*publicKeyECDSA.(*ecdsa.PublicKey))
	appstate.UpdateAccount(privateKeyHex, address.String())
	return nil
}

// LoadAccountFromPrivateKey loads an account from a private key.
func LoadAccountFromPrivateKey(privateKey string) error {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	appstate.UpdateAccount(privateKey, fromAddress.String())
	return nil
}

type BlockchainClient struct {
	client *ethclient.Client
}

func NewBlockchainClient(rpcURL string) (*BlockchainClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &BlockchainClient{client: client}, nil
}

// ERC20ABI is the standard ABI of the ERC20 token contract
const ERC20ABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

// GetERC20Balance gets the balance of a specific ERC20 token for an account.
func (bc *BlockchainClient) GetERC20Balance(tokenAddress string, accountAddress string) (*big.Int, error) {
	tokenAddressHex := common.HexToAddress(tokenAddress)
	accountAddressHex := common.HexToAddress(accountAddress)

	parsedABI, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}

	data, err := parsedABI.Pack("balanceOf", accountAddressHex)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddressHex,
		Data: data,
	}

	result, err := bc.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}
	var balanceResult interface{}
	err = parsedABI.UnpackIntoInterface(&balanceResult, "balanceOf", result)
	if err != nil {
		return nil, err
	}

	balance, ok := balanceResult.(*big.Int)
	if !ok {
		return nil, fmt.Errorf("balance result is not a big.Int")
	}

	return balance, nil
}

func (bc *BlockchainClient) GetETHBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := bc.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// CallReadOnlyFunction calls a read-only smart contract function.
func (bc *BlockchainClient) CallReadOnlyFunction(contractAddress string, abiJSON, functionName string, params ...interface{}) ([]interface{}, error) {
	// Parse the provided ABI
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, err
	}

	// Find the specific method in the ABI
	method, ok := parsedABI.Methods[functionName]
	if !ok {
		return nil, fmt.Errorf("function %s not found in ABI", functionName)
	}

	// Pack the parameters using the method's ABI
	packedParams, err := method.Inputs.Pack(params...)
	if err != nil {
		return nil, err
	}

	// Generate the function selector from the method ID
	selector := method.ID

	// Concatenate the selector and the packed parameters
	data := append(selector, packedParams...)

	address := common.HexToAddress(contractAddress)

	// Prepare the call message
	msg := ethereum.CallMsg{
		To:   &address,
		Data: data,
	}

	// Call the contract
	result, err := bc.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	// Unpack the output data
	output, err := method.Outputs.Unpack(result)
	if err != nil {
		return nil, err
	}

	return []interface{}{output[0]}, nil
}

// ExecuteWriteFunction executes a write operation on a smart contract function.
func (bc *BlockchainClient) ExecuteWriteFunction(privateKeyString, contractAddress, abiJSON, functionName string, gasLimit uint64, params ...interface{}) (common.Hash, error) {
	// Convert the private key string to an ecdsa.PrivateKey
	privateKeyBytes, err := hex.DecodeString(privateKeyString)
	if err != nil {
		return common.Hash{}, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return common.Hash{}, err
	}

	// Parse the provided ABI
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return common.Hash{}, err
	}

	// Find the specific method in the ABI
	method, ok := parsedABI.Methods[functionName]
	if !ok {
		return common.Hash{}, fmt.Errorf("function %s not found in ABI", functionName)
	}

	// Pack the parameters using the method's ABI
	packedParams, err := method.Inputs.Pack(params...)
	if err != nil {
		return common.Hash{}, err
	}

	// Generate the function selector from the method ID
	selector := method.ID

	// Concatenate the selector and the packed parameters
	data := append(selector, packedParams...)

	address := common.HexToAddress(contractAddress)

	// Fetch the nonce for the transaction
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := bc.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Hash{}, err
	}

	// Specify transaction parameters
	value := big.NewInt(0)
	gasPrice, err := bc.client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	// Create the transaction
	tx := types.NewTransaction(nonce, address, value, gasLimit, gasPrice, data)

	// Sign the transaction with the private key
	chainID, err := bc.client.NetworkID(context.Background())
	if err != nil {
		return common.Hash{}, err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	// Send the transaction
	err = bc.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, err
	}

	return signedTx.Hash(), nil
}

// WaitForTransactionReceipt waits for the transaction with the given hash to be mined.
func (bc *BlockchainClient) WaitForTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := bc.client.TransactionReceipt(context.Background(), txHash)
		if receipt != nil {
			return receipt, nil
		}
		if err != nil {
			time.Sleep(5 * time.Second) // wait for 5 seconds before trying again
		} else {
			return nil, err
		}
	}
}
