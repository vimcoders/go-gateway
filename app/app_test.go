package app

import (
	"context"
	"fmt"
	"net"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-lib"
	_ "github.com/vimcoders/sqlx-go-driver"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	sysLogger, err := lib.NewSyslogger()

	if err != nil {
		panic(err)
	}

	logger = sysLogger

	m.Run()
	fmt.Println("end")
}

func TestHandshake(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			c, err := net.Dial("tcp", ":8888")

			if err != nil {
				t.Error(err)
				return
			}

			s := Session{
				Conn:             c,
				v:                make(map[interface{}]interface{}),
				PushMessageQuene: make(chan driver.Message),
			}

			go s.Pull(context.Background())
		})
	}
}
