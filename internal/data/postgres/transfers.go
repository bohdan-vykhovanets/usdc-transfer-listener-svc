package postgres

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/logan/v3/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const tableName = "transfers"

type transferQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

func newTransferQ(db *pgdb.DB) data.TransferQ {
	baseSelect := sq.Select("*").From(tableName)

	return &transferQ{
		db:  db,
		sql: baseSelect,
	}
}

func (q *transferQ) Select() (*[]data.Transfer, error) {
	var result []data.Transfer

	err := q.db.Select(&result, q.sql)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select transfers from database")
	}
	if result == nil {
		result = []data.Transfer{}
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

func (q *transferQ) FilterByFrom(from string) data.TransferQ {
	pred := sq.Eq{`"from"`: from}
	q.sql = q.sql.Where(pred)
	return q
}

func (q *transferQ) FilterByTo(to string) data.TransferQ {
	pred := sq.Eq{`"to"`: to}
	q.sql = q.sql.Where(pred)
	return q
}

func (q *transferQ) FilterByCounterparty(counterparty string) data.TransferQ {
	pred := sq.Or{
		sq.Eq{`"from"`: counterparty},
		sq.Eq{`"to"`: counterparty},
	}
	q.sql = q.sql.Where(pred)
	return q
}

func (q *transferQ) Paginate(limit, offset uint64) data.TransferQ {
	q.sql = q.sql.Limit(limit).Offset(offset)
	return q
}
