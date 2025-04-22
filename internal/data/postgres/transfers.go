package postgres

import (
	sql2 "database/sql"
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
		Columns(`"block_number"`, `"tx_hash"`, `"log_index"`, `"from"`, `"to"`, `"value"`).
		Values(
			value.BlockNumber,
			value.TxHash,
			value.LogIndex,
			value.From,
			value.To,
			value.Value.Int.String(),
		).
		Suffix("ON CONFLICT (block_number, tx_hash, log_index) DO NOTHING")
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

func (q *transferQ) GetLastProcessedBlock() (uint64, error) {
	var result sql2.NullInt64

	sql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	stmt := sql.Select("MAX(block_number)").From(tableName)
	err := q.db.Get(&result, stmt)
	if err != nil {
		return 0, errors.Wrap(err, "failed to query last processed block")
	}

	if !result.Valid {
		return 0, nil
	}

	return uint64(result.Int64), nil
}
