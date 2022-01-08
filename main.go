package main

import (
	"github.com/vimcoders/go-gateway/ctx"
	"github.com/vimcoders/go-gateway/lib"
	"github.com/vimcoders/go-gateway/logx"
	"github.com/vimcoders/go-gateway/sqlx"

	_ "github.com/vimcoders/go-gateway/mongox"
	_ "github.com/vimcoders/go-gateway/session"
)

func init() {
	logx.Info("init cost %vs", lib.Seconds())
}

func main() {
	closeCtx := ctx.Close()

	for {
		select {
		case <-closeCtx.Done():
			sqlx.Close()
			logx.Info("shutdown %vs", lib.Seconds())
			return
		}
	}
}
