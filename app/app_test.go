package app

import (
	"context"
	"fmt"
	"net"
	"sync"
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
	var waitGroup sync.WaitGroup

	for i := 0; i < 10; i++ {
		waitGroup.Add(1)

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
				OnMessage: func(pkg driver.Message) (err error) {
					b, err := pkg.ToBytes()

					if err != nil {
						logger.Error("OnMessage %v", err)
						return err
					}

					logger.Info("OnMessage %v", string(b))

					waitGroup.Done()

					return nil
				},
			}

			go s.Pull(context.Background())
		})
	}

	waitGroup.Wait()
}
