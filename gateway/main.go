package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/vimcoders/go-gateway/log"
	"github.com/vimcoders/go-gateway/session"
)

var (
	addr                = ":8888"
	httpAddr            = "localhost:8000"
	network             = "tcp"
	closeCtx, closeFunc = context.WithCancel(context.Background())
)

func main() {
	Run()
}

func Run() (err error) {
	now := time.Now()

	defer func() {
		if e := recover(); e != nil {
			log.Error("Run Recover %v", e)
		}

		if err != nil {
			log.Error("Run Recover %v", err)
		}

		closeFunc()
	}()

	var waitGroup sync.WaitGroup

	waitGroup.Add(2)

	go Listen(&waitGroup)

	go Monitor(&waitGroup)

	log.Info("Run Cost %v", time.Now().Sub(now))

	waitGroup.Wait()

	return errors.New("shutdown!")
}

func Listen(waitGroup *sync.WaitGroup) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("Listen %v", e)
		}

		if err != nil {
			log.Error("Listen %v", err)
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
				log.Error("Listen %v", err.Error())
				continue
			}

			go session.Handle(closeCtx, conn)
		}
	}
}

func Monitor(waitGroup *sync.WaitGroup) (err error) {
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
