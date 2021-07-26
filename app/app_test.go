package app

import (
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
	m.Run()
	fmt.Println("end")
}

func TestLogin(t *testing.T) {
	type Client struct {
		net.Conn
		OnMessage        func(pkg driver.Message) (err error)
		sender           *Encoder
		reader           *Reader
		PushMessageQuene chan driver.Message
	}

	var waitGroup sync.WaitGroup

	for i := 0; i < 2; i++ {
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
			client.reader = NewReader(c, NewBuffer())
			client.PushMessageQuene = make(chan driver.Message, 1)
			client.OnMessage = func(pkg driver.Message) (err error) {
				b, err := pkg.ToBytes()

				if err != nil {
					return err
				}

				t.Logf("OnMessage %v", string(b))

				block, result := pem.Decode(b)

				if len(result) > 0 {
					if err = client.sender.Write(pkg); err != nil {
						return
					}

					return nil
				}

				key, err := x509.ParsePKIXPublicKey(block.Bytes)

				if err != nil {
					return err
				}

				publicKey := key.(*rsa.PublicKey)
				client.sender = NewEncoder(c, NewBuffer(), publicKey)

				if err = client.sender.Write(NewMessage([]byte("hello server !"))); err != nil {
					return
				}

				return nil
			}

			defer func() {
				if err := client.Close(); err != nil {
					t.Errorf("close err %v", err)
				}
			}()

			for {
				pkg, err := client.reader.Read()

				if err != nil {
					t.Errorf("OnMessage %v", err)
					return
				}

				if err := client.OnMessage(pkg); err != nil {
					t.Errorf("OnMessage %v", err)
					return
				}

				time.Sleep(time.Second)
			}
		})
	}

	waitGroup.Wait()
}
