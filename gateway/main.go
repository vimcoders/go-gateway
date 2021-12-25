package main

import (
	"errors"
	"net/http"
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

	mongox.Init()
	sqlx.Init()
	session.Init()
	go http.ListenAndServe(":8080", nil)

	log.Info("Run Cost %v", time.Now().Sub(now))

	ctx := lib.Context()

	for {
		select {
		case <-ctx.Done():
			return errors.New("shutdown")
		}
	}
}
