package apk

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/vimcoders/go-driver"

	"net/http"
	_ "net/http/pprof"
)

var (
	addr     = ":8888"
	httpAddr = "localhost:8000"
	network  = "tcp"

	moment              = time.Now()
	logger, _           = driver.NewSyslogger()
	closeCtx, closeFunc = context.WithCancel(context.Background())
)

func Listen(waitGroup *sync.WaitGroup) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("Listen %v", e)
		}

		if err != nil {
			logger.Error("Listen %v", err)
		}

		waitGroup.Done()
	}()

	listener, err := net.Listen(network, addr)

	if err != nil {
		return err
	}

	for {
		select {
		case <-closeCtx.Done():
			return errors.New("shutdown")
		default:
			conn, err := listener.Accept()

			if err != nil {
				logger.Error("Listen %v", err.Error())
				continue
			}

			go Handle(closeCtx, conn)
		}
	}
}

func Monitor(waitGroup *sync.WaitGroup) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("Listen %v", e)
		}

		if err != nil {
			logger.Error("Listen %v", err)
		}

		waitGroup.Done()
	}()

	return http.ListenAndServe(httpAddr, nil)
}

func Run() (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("Run Recover %v", e)
		}

		if err != nil {
			logger.Error("Run Recover %v", err)
		}

		closeFunc()
	}()

	var waitGroup sync.WaitGroup

	waitGroup.Add(2)

	go Listen(&waitGroup)

	go Monitor(&waitGroup)

	logger.Info("Run Cost %v", time.Now().Sub(moment))

	waitGroup.Wait()

	return errors.New("shutdown!")
}
