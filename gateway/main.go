package main

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/vimcoders/go-gateway/lib"
	"github.com/vimcoders/go-gateway/log"
	"github.com/vimcoders/go-gateway/mongox"
	"github.com/vimcoders/go-gateway/session"
	"github.com/vimcoders/go-gateway/sqlx"
)

func main() {
	run()
}

func run() (err error) {
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

	wg.Add(4)
	go monitor(&wg)
	mongox.Init(&wg)
	sqlx.Init(&wg)
	session.Init(&wg)

	log.Info("Run Cost %v", time.Now().Sub(now))

	wg.Wait()

	return errors.New("shutdown!")
}

func monitor(wg *sync.WaitGroup) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("Listen %v", e)
		}

		if err != nil {
			log.Error("Listen %v", err)
		}

		wg.Done()
	}()

	addr := lib.MonitorAddr()

	return http.ListenAndServe(addr, nil)
}
