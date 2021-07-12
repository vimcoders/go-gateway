package app

import (
	"context"
	"net"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-lib"
	"github.com/vimcoders/sqlx-go-driver"
)

var logger driver.Logger
var connector driver.Connector

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

//func Monitor(waitGroup *sync.WaitGroup) (err error) {
//	defer func() {
//		if e := recover(); e != nil {
//			logger.Error("Listen %v", e)
//		}
//
//		if err != nil {
//			logger.Error("Listen %v", err)
//		}
//
//		waitGroup.Done()
//	}()
//
//	http.Handle("/metrics", promhttp.Handler())
//	return http.ListenAndServe(":2112", nil)
//}

func Run() {
	now := time.Now()

	sysLogger, err := lib.NewSyslogger()

	if err != nil {
		panic(err)
	}

	logger = sysLogger

	sqlConnector, err := sqlx.Connect(&sqlx.Config{
		DriverName: "mysql",
	})

	if err != nil {
		panic(err)
	}

	connector = sqlConnector

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go Listen(&waitGroup)

	//go Monitor(&waitGroup)

	logger.Info("Run Cost %v", time.Now().Sub(now))

	waitGroup.Wait()
}
