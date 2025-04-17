package postgres

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/logan/v3/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const tableName = "transfers"

func newTransferQ(db *pgdb.DB) data.TransferQ {
	return &transferQ{
		db:  db,
		sql: sq.StatementBuilder,
	}
}

type transferQ struct {
	db  *pgdb.DB
	sql sq.StatementBuilderType
}

func (q *transferQ) Select() (*[]data.Transfer, error) {
	var result []data.Transfer

	stmt := sq.Select(`"id", "from", "to", "value"`).From(tableName)
	err := q.db.Select(&result, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select transfers from database")
	}

	return &result, nil
}

func (q *transferQ) Insert(value data.Transfer) error {
	stmt := sq.Insert(tableName).
		Columns(`"from"`, `"to"`, `"value"`).
		Values(
			common.Address(value.From).Hex(),
			common.Address(value.To).Hex(),
			value.Value.Int.String(),
		)
	err := q.db.Exec(stmt)
	if err != nil {
		return errors.Wrap(err, "failed to insert transfer to database")
	}

	return nil
}
