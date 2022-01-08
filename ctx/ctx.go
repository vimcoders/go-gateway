package ctx

import "context"

var (
	closeCtx, closeFunc = context.WithCancel(context.Background())
)

func Close() context.Context {
	return closeCtx
}
