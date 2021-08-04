package app

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/vimcoders/go-driver"

	"net/http"
	_ "net/http/pprof"
)

var (
	logger     driver.Logger
	closeFunc  context.CancelFunc
	closeCtx   context.Context
	addr       = ":8888"
	httpAddr   = "localhost:8000"
	network    = "tcp"
	privateKey *rsa.PrivateKey
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

			go Handle(closeCtx, conn, privateKey)
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
	closeCtx, closeFunc = context.WithCancel(context.Background())
	defer closeFunc()

	now := time.Now()

	sysLogger, err := driver.NewSyslogger()

	if err != nil {
		panic(err)
	}

	logger = sysLogger

	defer func() {
		if e := recover(); e != nil {
			logger.Error("Run Recover %v", e)
		}

		if err != nil {
			logger.Error("Run Recover %v", err)
		}
	}()

	key, err := rsa.GenerateKey(rand.Reader, 512)

	if err != nil {
		return err
	}

	privateKey = key

	var waitGroup sync.WaitGroup

	waitGroup.Add(2)

	go Listen(&waitGroup)

	go Monitor(&waitGroup)

	logger.Info("Run Cost %v", time.Now().Sub(now))

	waitGroup.Wait()

	return errors.New("shutdown!")
}
