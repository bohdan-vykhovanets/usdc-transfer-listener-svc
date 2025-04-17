package background

import (
	"context"
	"fmt"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data/dbtypes"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/token"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"sync"
)

func Transfers(ctx context.Context, client *ethclient.Client, wg *sync.WaitGroup) {
	Log(ctx).Info("In Transfers function...")
	defer wg.Done()

	var TransferEventSignature = []byte("Transfer(address,address,uint256)")
	const USDCTokenAddress = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(USDCTokenAddress)},
		Topics: [][]common.Hash{
			{crypto.Keccak256Hash(TransferEventSignature)},
		},
	}
	Log(ctx).Info("Query created: ", query)

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		Log(ctx).Errorf("failed to subscribe to logs: %v", err)
		return
	}
	defer sub.Unsubscribe()
	Log(ctx).Info("Subscribed to logs")

	for {
		Log(ctx).Info("Reading logs...")
		select {
		case err := <-sub.Err():
			Log(ctx).Errorf("log subscription error: %v", err)
			return
		case log := <-logs:
			parsedLog, err := parseTransferLog(log)
			if err != nil {
				Log(ctx).Errorf("failed to parse log: %v", err)
			}
			err = Db(ctx).Transfer().Insert(*parsedLog)
			if err != nil {
				Log(ctx).Errorf("failed to insert transfer: %v", err)
			}
			Log(ctx).Info("Parsed logs")
		case <-ctx.Done():
			return
		}
	}
}

func parseTransferLog(log types.Log) (*data.Transfer, error) {
	contractAbi, err := token.NewERC20Filterer(log.Address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate filterer: %v", err)
	}

	transfer, err := contractAbi.ParseTransfer(log)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log as transfer: %v", err)
	}

	return &data.Transfer{
		From: dbtypes.DbAddress(transfer.From),
		To:   dbtypes.DbAddress(transfer.To),
		Value: dbtypes.DbBigInt{
			Int: transfer.Value,
		},
	}, nil
}
