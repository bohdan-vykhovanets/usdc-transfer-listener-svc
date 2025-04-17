package background

import (
	"context"
	"github.com/bohdan-vykhovanets/usdc-transfer-listener-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	dbCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func CtxDb(entry data.MainQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context { return context.WithValue(ctx, dbCtxKey, entry) }
}

func Log(ctx context.Context) *logan.Entry {
	return ctx.Value(logCtxKey).(*logan.Entry)
}

func Db(ctx context.Context) data.MainQ {
	return ctx.Value(dbCtxKey).(data.MainQ).New()
}
