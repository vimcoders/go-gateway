package main

import (
	"github.com/vimcoders/go-gateway/lib"
	"github.com/vimcoders/go-gateway/log"

	_ "github.com/vimcoders/go-gateway/mongox"
	_ "github.com/vimcoders/go-gateway/session"
	_ "github.com/vimcoders/go-gateway/sqlx"
)

func init() {
	log.Info("init cost %v", lib.Duration())
}

func main() {
	ctx := lib.Context()

	for {
		select {
		case <-ctx.Done():
		}
	}
}
