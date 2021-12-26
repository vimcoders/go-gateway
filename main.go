package main

import (
	"github.com/vimcoders/go-gateway/lib"
	"github.com/vimcoders/go-gateway/log"
	"github.com/vimcoders/go-gateway/sqlx"

	_ "github.com/vimcoders/go-gateway/mongox"
	_ "github.com/vimcoders/go-gateway/session"
	_ "github.com/vimcoders/go-gateway/sqlx"
)

func init() {
	log.Info("init cost %vs", lib.Seconds())
}

func main() {
	ctx := lib.Context()

	for {
		select {
		case <-ctx.Done():
			sqlx.Close()
			log.Info("shutdown %vs", lib.Seconds())
			return
		}
	}
}
