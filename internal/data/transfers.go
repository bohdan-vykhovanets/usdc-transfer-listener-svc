package data

import (
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data/dbtypes"
)

type TransferQ interface {
	Select() (*[]Transfer, error)
	Insert(value Transfer) error
}

type Transfer struct {
	ID    int64             `db:"id" structs:"-"`
	From  dbtypes.DbAddress `db:"from" structs:"from"`
	To    dbtypes.DbAddress `db:"to" structs:"to"`
	Value dbtypes.DbBigInt  `db:"value" structs:"value"`
}
