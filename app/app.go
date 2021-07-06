package app

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-lib"
	"github.com/vimcoders/pb"
)

var logger driver.Logger
var grpc pb.LogicClient

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
			logger.Error("Listen %v", err.Error())
			continue
		}

		Handle(context.Background(), conn)
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

	http.Handle("/metrics", promhttp.Handler())
	return	http.ListenAndServe(":2112", nil)
}

func Run() {
	now := time.Now()

	sysLogger, err := lib.NewSyslogger()

	if err != nil {
		panic(err)
	}

	logger = sysLogger

	var waitGroup sync.WaitGroup

	waitGroup.Add(2)

	go Listen(&waitGroup)

	go Monitor(&waitGroup)

	logger.Info("Run Cost %v", time.Now().Sub(now))

	waitGroup.Wait()
}
