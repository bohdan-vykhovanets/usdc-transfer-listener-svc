package data

import (
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data/dbtypes"
)

type TransferQ interface {
	Select() (*[]Transfer, error)
	Insert(value Transfer) error

	FilterByFrom(from string) TransferQ
	FilterByTo(to string) TransferQ
	FilterByCounterparty(counterparty string) TransferQ
	Paginate(limit, offset uint64) TransferQ

	GetLastProcessedBlock() (uint64, error)
}

type Transfer struct {
	ID          int64             `db:"id" structs:"-"`
	BlockNumber uint64            `db:"block_number" structs:"block_number"`
	TxHash      dbtypes.DbHash    `db:"tx_hash" structs:"tx_hash"`
	LogIndex    uint              `db:"log_index" structs:"log_index"`
	From        dbtypes.DbAddress `db:"from" structs:"from"`
	To          dbtypes.DbAddress `db:"to" structs:"to"`
	Value       dbtypes.DbBigInt  `db:"value" structs:"value"`
}
