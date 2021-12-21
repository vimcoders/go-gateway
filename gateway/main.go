package main

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/vimcoders/go-gateway/log"
	//_ "github.com/vimcoders/go-gateway/mongo"
	_ "github.com/vimcoders/go-gateway/session"
	//_ "github.com/vimcoders/go-gateway/sqlx"
)

var (
	httpAddr = "localhost:8000"
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
	}()

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go Monitor(&waitGroup)

	log.Info("Run Cost %v", time.Now().Sub(now))

	waitGroup.Wait()

	return errors.New("shutdown!")
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
