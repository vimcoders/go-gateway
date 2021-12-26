package lib

import (
	"context"
	"net/http"
	"time"
)

var (
	addr                = ":8888"
	unix                = time.Now().Unix()
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

func Seconds() int64 {
	return time.Now().Unix() - unix
}

func Shutdown() {
	closeFunc()
}
