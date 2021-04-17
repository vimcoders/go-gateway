package app

import (
	"context"
	"net"
	"sync"

	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-lib"
)

var logger driver.Logger

func Listen() (err error) {
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
	sysLogger, err := lib.NewSyslogger()

	if err != nil {
		panic(err)
	}

	logger = sysLogger

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
}
