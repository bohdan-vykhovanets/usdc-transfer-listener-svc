package background

import (
	"context"
	"fmt"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data/dbtypes"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"sync"
)

const (
	USDCTokenAddress  = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
	defaultStartBlock = 22326050
	defaultBatchSize  = 100
)

func Transfers(ctx context.Context, client *ethclient.Client, wg *sync.WaitGroup) {
	logger := Log(ctx)

	logger.Info("Started transfers background job...")
	defer wg.Done()

	err := syncPastTransfers(ctx, client)
	if err != nil {
		logger.Errorf("Failed to sync past transfers: %v", err)
		return
	}

	usdcAddress := common.HexToAddress(USDCTokenAddress)
	erc20Filterer, err := token.NewERC20Filterer(usdcAddress, client)
	if err != nil {
		logger.Errorf("Failed to create ERC20Filterer for sunbscription: %v", err)
		return
	}

	events := make(chan *token.ERC20Transfer)

	select {
	case <-ctx.Done():
		return
	default:
	}

	watchOptions := &bind.WatchOpts{
		Context: ctx,
	}

	sub, err := erc20Filterer.WatchTransfer(watchOptions, events, nil, nil)
	if err != nil {
		logger.Errorf("failed to subscribe to logs: %v", err)
		return
	}
	defer sub.Unsubscribe()
	logger.Info("Subscribed to logs")

	for {
		select {
		case err := <-sub.Err():
			logger.Errorf("log subscription error: %v", err)
			return
		case event := <-events:
			if event.Raw.Removed {
				continue
			}

			parsedTransfer := &data.Transfer{
				BlockNumber: event.Raw.BlockNumber,
				TxHash:      dbtypes.DbHash(event.Raw.TxHash),
				LogIndex:    event.Raw.Index,
				From:        dbtypes.DbAddress(event.From),
				To:          dbtypes.DbAddress(event.To),
				Value: dbtypes.DbBigInt{
					Int: event.Value,
				},
			}

			err = Db(ctx).Transfer().Insert(*parsedTransfer)
			if err != nil {
				logger.Errorf("failed to insert transfer: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func syncPastTransfers(ctx context.Context, client *ethclient.Client) error {
	logger := Log(ctx)

	logger.Info("Started transfers synchronization...")
	lastBlock, err := Db(ctx).Transfer().GetLastProcessedBlock()
	if err != nil {
		logger.Errorf("failed to get last processed block: %v", err)
		return fmt.Errorf("failed to get last processed block: %v", err)
	}

	startBlock := lastBlock + 1
	if lastBlock == 0 {
		startBlock = defaultStartBlock
		logger.Infof("Starting with default block %d", startBlock)
	} else {
		logger.Infof("Starting with last block +1 %d", startBlock)
	}

	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		logger.Errorf("failed to get latest block: %v", err)
		return fmt.Errorf("failed to get latest block: %v", err)
	}

	usdcAddress := common.HexToAddress(USDCTokenAddress)
	erc20Filterer, err := token.NewERC20Filterer(usdcAddress, client)
	if err != nil {
		return fmt.Errorf("failed to instantiate ERC20Filterer: %w", err)
	}

	for fromBlock := startBlock; fromBlock <= latestBlock; fromBlock += defaultBatchSize {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		toBlock := fromBlock + defaultBatchSize - 1
		if toBlock > latestBlock {
			toBlock = latestBlock
		}

		logger.Infof("Fetching blocks from %d to %d", fromBlock, toBlock)

		filterOptions := &bind.FilterOpts{
			Start:   fromBlock,
			End:     &toBlock,
			Context: ctx,
		}

		iterator, err := erc20Filterer.FilterTransfer(filterOptions, nil, nil)
		if err != nil {
			logger.Errorf("failed to filter logs: %v", err)
			return fmt.Errorf("failed to filter logs: %v", err)
		}

		for iterator.Next() {
			event := iterator.Event
			if event.Raw.Removed {
				continue
			}

			parsedTransfer := &data.Transfer{
				BlockNumber: event.Raw.BlockNumber,
				TxHash:      dbtypes.DbHash(event.Raw.TxHash),
				LogIndex:    event.Raw.Index,
				From:        dbtypes.DbAddress(event.From),
				To:          dbtypes.DbAddress(event.To),
				Value: dbtypes.DbBigInt{
					Int: event.Value,
				},
			}

			err = Db(ctx).Transfer().Insert(*parsedTransfer)
			if err != nil {
				logger.Errorf("failed to insert transfer: %v", err)
			}
		}

		if err := iterator.Error(); err != nil {
			logger.Errorf("failed to iterate logs: %v", err)
			iterator.Close()
			return fmt.Errorf("failed to iterate logs: %v", err)
		}

		if err := iterator.Close(); err != nil {
			logger.Errorf("failed to close iterator: %v", err)
		}
	}

	logger.Info("Finished transfers synchronization")

	return nil
}
