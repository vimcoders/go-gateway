package sqlx

import (
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-gateway/log"
	"github.com/vimcoders/sqlx-go-driver"
)

var connector driver.Connector

func init() {
	c, err := sqlx.Connect(&sqlx.Config{
		DriverName: "mysql",
		Usr:        "centos",
		Pwd:        "cenots",
		Addr:       "localhost",
	})

	if err != nil {
		log.Error("err %v", err)
		os.Exit(0)
		return
	}

	connector = c
	log.Info("mysql init success......")
}

func Close() error {
	if err := connector.Close(); err != nil {
		log.Error("sqlx close err %v", err.Error())
		return err
	}

	return nil
}
