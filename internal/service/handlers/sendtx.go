package handlers

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/ape"
	"log"
	"math/big"
	"net/http"
	"os"
)

func SendTransactionTest(w http.ResponseWriter, r *http.Request) {
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

	nonce, err := client.PendingNonceAt(context.Background(), testAddress)
	if err != nil {
		log.Fatalf("failed to get nonce: %v", err)
	}
	logger.Infof("current nonce: %d", nonce)

	value := big.NewInt(10e15)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to get gas price: %v", err)
	}
	gasLimit := uint64(21000)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &testAddress,
		Value:    value,
		Data:     nil,
	})

	ChainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to get network ID: %v", err)
	}
	signer := types.NewEIP155Signer(ChainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		logger.Fatalf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		logger.Fatalf("failed to send transaction: %v", err)
	}
	logger.Infof("transaction sent")

	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		logger.Fatalf("failed to wait for receipt: %v", err)
	}
	logger.Infof("transaction mined, status: %v", receipt.Status)

	ape.Render(w, receipt.Status)
}
