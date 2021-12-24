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

	var wg sync.WaitGroup

	wg.Add(3)
	mongox.Init(&wg)
	sqlx.Init(&wg)
	session.Init(&wg)

	log.Info("Run Cost %v", time.Now().Sub(now))

	wg.Wait()

	return errors.New("shutdown!")
}
