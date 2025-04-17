package postgres

import (
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func NewMainQ(db *pgdb.DB) data.MainQ {
	return &mainQ{
		db: db.Clone(),
	}
}

type mainQ struct {
	db *pgdb.DB
}

func (m *mainQ) New() data.MainQ {
	return NewMainQ(m.db)
}

func (m *mainQ) Transfer() data.TransferQ {
	return newTransferQ(m.db)
}
