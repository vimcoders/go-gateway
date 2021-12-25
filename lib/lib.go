package lib

import (
	"context"
	"net/http"
	"time"
)

var (
	addr                = ":8888"
	awake               = time.Now()
	timeout             = time.Second * 15
	monitorAddr         = ":8080"
	closeCtx, closeFunc = context.WithCancel(context.Background())
)

func init() {
	go http.ListenAndServe(":8080", nil)
}

func Context() context.Context {
	return closeCtx
}

func Timeout() time.Time {
	return time.Now().Add(timeout)
}

func Addr() string {
	return addr
}

func Duration() time.Duration {
	return time.Now().Sub(awake)
}

func Shutdown() {
	closeFunc()
}
