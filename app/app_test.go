package app

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/vimcoders/go-driver"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	m.Run()
	fmt.Println("end")
}

func TestLogin(t *testing.T) {
	type Client struct {
		Session
	}

	var waitGroup sync.WaitGroup

	for i := 0; i < 20000000; i++ {
		waitGroup.Add(1)

		go t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			defer waitGroup.Done()

			c, err := net.Dial("tcp", ":8888")

			if err != nil {
				t.Error(err)
				return
			}

			var client Client

			client.Conn = c
			client.Buffer = NewBuffer()
			client.PushMessageQuene = make(chan driver.Message, 1)
			client.OnMessage = func(pkg driver.Message) (err error) {
				b, err := pkg.ToBytes()

				if err != nil {
					return err
				}

				t.Logf("rsa key %v", string(b))

				block, result := pem.Decode(b)

				if len(result) > 0 {
					return errors.New(fmt.Sprintf("pem result %v", string(result)))
				}

				key, err := x509.ParsePKIXPublicKey(block.Bytes)

				if err != nil {
					return err
				}

				publicKey := key.(*rsa.PublicKey)

				coder := NewEncoder(publicKey, []byte("hello golang"))

				if err := client.Send(coder); err != nil {
					return err
				}

				return nil
			}

			if err := client.WaitMessage(); err != nil {
				t.Errorf("WaitMessage %v", err)
			}

			client.Close()
		})
	}

	waitGroup.Wait()
}
