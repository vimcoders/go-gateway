package app

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-lib"
)

var logger driver.Logger

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

	listener, err := net.Listen("tcp", ":8888")

	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		NewSession(context.Background(), conn)
	}
}

func Run() {
	now := time.Now()

	sysLogger, err := lib.NewSyslogger()

	if err != nil {
		panic(err)
	}

	logger = sysLogger

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go Listen(&waitGroup)

	logger.Info("Run Cost %v", time.Now().Sub(now))

	waitGroup.Wait()
}
