package handlers

import (
	"context"
	"crypto/ecdsa"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/ape"
	"log"
	"math/big"
	"net/http"
	"os"
)

const tokenAddressHex = "0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238" // USDC address in Sepolia, 6 decimals

func SendTokenTest(w http.ResponseWriter, r *http.Request) {
	logger := Log(r)
	nodeUrl := Node(r)
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if nodeUrl == "" {
		logger.Fatal("node URL is empty")
	}
	if privateKeyHex == "" {
		logger.Fatal("private key is empty")
	}

	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		logger.Fatalf("failed to connect to node: %v", err)
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		logger.Fatalf("failed to parse private key: %v", err)
	}
	publicKeyRaw := privateKey.Public()
	publicKey, ok := publicKeyRaw.(*ecdsa.PublicKey)
	if !ok {
		logger.Fatal("error casting public key")
	}

	testAddress := crypto.PubkeyToAddress(*publicKey)
	logger.Infof("testing address: %s", testAddress.Hex())
	tokenAddress := common.HexToAddress(tokenAddressHex)
	logger.Infof("testing token: %s", tokenAddress.Hex())

	nonce, err := client.PendingNonceAt(context.Background(), testAddress)
	if err != nil {
		log.Fatalf("failed to get nonce: %v", err)
	}
	logger.Infof("current nonce: %d", nonce)

	value := big.NewInt(10e6) // = 1 USDC

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to get gas price: %v", err)
	}

	ChainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to get network ID: %v", err)
	}
	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, ChainID)
	if err != nil {
		logger.Fatalf("failed to create keyed transactor: %v", err)
	}
	transactor.Nonce = big.NewInt(int64(nonce))
	transactor.Value = big.NewInt(0)
	transactor.GasLimit = uint64(0)
	transactor.GasPrice = gasPrice

	tokenInst, err := token.NewErc20(tokenAddress, client)
	if err != nil {
		logger.Fatalf("failed to create token instance: %v", err)
	}

	tx, err := tokenInst.Transfer(transactor, testAddress, value)
	if err != nil {
		logger.Fatalf("failed to send token transaction: %v", err)
	}

	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		logger.Fatalf("failed to wait for receipt: %v", err)
	}
	logger.Infof("transaction mined, status: %v", receipt.Status)

	ape.Render(w, receipt.Status)
}
