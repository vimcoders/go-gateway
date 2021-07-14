package app

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-gateway/pb"
	"github.com/vimcoders/go-lib"
	_ "github.com/vimcoders/sqlx-go-driver"
	"google.golang.org/protobuf/proto"
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

func TestLogin(t *testing.T) {
	var waitGroup sync.WaitGroup

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)

		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			defer waitGroup.Done()

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

			s.OnMessage = func(pkg driver.Message) (err error) {
				b, err := pkg.ToBytes()

				if err != nil {
					logger.Error("OnMessage %v", err)
					return err
				}

				logger.Info("%v", string(b))

				block, _ := pem.Decode(b)

				key, err := x509.ParsePKIXPublicKey(block.Bytes)

				if err != nil {
					logger.Error("OnMessage %v", err)
					return err
				}

				publicKey := key.(*rsa.PublicKey)

				login := &pb.Login{UserName: "golangxxxxxxxxxx", Pwd: "golangxxxxxxxxxx"}

				loginBytes, err := proto.Marshal(login)

				if err != nil {
					logger.Error("marshal %v", err)
					return err
				}

				coder := NewEncoder(publicKey, loginBytes)

				if err := s.Write(coder); err != nil {
					logger.Error("encoder %v", err)
					return err
				}

				return nil
			}

			go s.Push(context.Background())
			s.Pull(context.Background())
		})
	}

	waitGroup.Wait()
}
