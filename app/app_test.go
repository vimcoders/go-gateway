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
	"github.com/vimcoders/go-lib"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	m.Run()
	fmt.Println("end")
}

func TestLogin(t *testing.T) {
	type Client struct {
		key *rsa.PublicKey
		net.Conn
		driver.Writer
		driver.Reader
		OnMessage        func(pkg driver.Message) (err error)
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
			client.Reader = lib.NewReader(c, lib.NewBuffer(), time.Second*5)
			client.Writer = lib.NewWriter(c, lib.NewBuffer(), time.Second*5)
			client.PushMessageQuene = make(chan driver.Message, 1)
			client.OnMessage = func(pkg driver.Message) (err error) {
				b, err := pkg.ToBytes()

				if err != nil {
					return err
				}

				t.Logf("OnMessage %v", string(b))

				block, _ := pem.Decode(b)
				key, err := x509.ParsePKIXPublicKey(block.Bytes)

				if err != nil {
					return err
				}

				client.key = key.(*rsa.PublicKey)

				encoder := lib.NewEncoder([]byte("hello server !"), client.key)

				if err = client.Writer.Write(encoder); err != nil {
					return err
				}

				return nil
			}

			defer func() {
				if err := client.Close(); err != nil {
					t.Errorf("close err %v", err)
				}
			}()

			for {
				pkg, err := client.Reader.Read()

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
