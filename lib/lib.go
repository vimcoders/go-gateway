package lib

import (
	"context"
	"time"
)

var (
	addr                = ":8888"
	monitorAddr         = ":8080"
	timeout             = time.Second * 15
	closeCtx, closeFunc = context.WithCancel(context.Background())
)

func Context() context.Context {
	return closeCtx
}

func Timeout() time.Time {
	return time.Now().Add(timeout)
}

func Addr() string {
	return addr
}

func MonitorAddr() string {
	return monitorAddr
}
