package main

import (
	"errors"
	"sync"
	"time"

	"github.com/vimcoders/go-gateway/log"
	"github.com/vimcoders/go-gateway/mongox"
	"github.com/vimcoders/go-gateway/session"
	"github.com/vimcoders/go-gateway/sqlx"
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

	waitGroup.Add(3)
	mongox.Init(&waitGroup)
	sqlx.Init(&waitGroup)
	session.Init(&waitGroup)

	log.Info("Run Cost %v", time.Now().Sub(now))

	waitGroup.Wait()

	return errors.New("shutdown!")
}
