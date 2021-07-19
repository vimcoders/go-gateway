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
	"time"

	"github.com/vimcoders/go-driver"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	go Run()

	m.Run()
	fmt.Println("end")
}

func TestLogin(t *testing.T) {
	var waitGroup sync.WaitGroup

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		time.Sleep(time.Second)

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

				coder := NewEncoder(publicKey, []byte("hello golang"))

				if err := s.Send(coder); err != nil {
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
