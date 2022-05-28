package main

import (
	"context"
	"net/http"

	"github.com/vimcoders/go-gateway/logx"
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
			logx.Info("shutdown")
			return
		}
	}
}
