package main

import (
	"context"
	"net/http"

	"github.com/vimcoders/go-gateway/logx"
	"github.com/vimcoders/go-gateway/session"
	"github.com/vimcoders/go-gateway/sqlx"
)

var (
	CloseCtx, CloseFunc = context.WithCancel(context.Background())
)

func init() {
	go http.ListenAndServe(":8080", nil)
}

func main() {
	for {
		select {
		case <-CloseCtx.Done():
			sqlx.Close()
			session.CloseFunc()
			logx.Info("shutdown")
			return
		}
	}
}
