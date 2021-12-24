package lib

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/vimcoders/go-gateway/log"
)

var (
	addr                = ":8888"
	closeCtx, closeFunc = context.WithCancel(context.Background())
	timeout             = time.Second * 15
	httpAddr            = "localhost:8000"
)

func Init(wg *sync.WaitGroup) {
	go monitor(wg)
}

func monitor(waitGroup *sync.WaitGroup) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("Listen %v", e)
		}

		if err != nil {
			log.Error("Listen %v", err)
		}

		waitGroup.Done()
	}()

	return http.ListenAndServe(httpAddr, nil)
}

func Context() (context.Context, context.CancelFunc) {
	return closeCtx, closeFunc
}

func Timeout() time.Time {
	return time.Now().Add(timeout)
}

func Addr() string {
	return addr
}
