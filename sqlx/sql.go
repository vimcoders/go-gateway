package sqlx

import (
	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-gateway/log"
	sqlx "github.com/vimcoders/sqlx-go-driver"
)

var connector driver.Connector

func init() {
	c, err := sqlx.Connect(nil)

	if err != nil {
		log.Error("err %v", err)
	}

	connector = c
}
