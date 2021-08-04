package app

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	m.Run()
	fmt.Println("end")
}

func TestLogin(t *testing.T) {
	type Client struct {
		key *rsa.PublicKey
		io.Closer
		io.Writer
		*Reader
		OnMessage func(pkg []byte) (err error)
	}

	var waitGroup sync.WaitGroup

	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)

		go t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			defer waitGroup.Done()

			c, err := net.Dial("tcp", ":8888")

			if err != nil {
				t.Error(err)
				return
			}

			var client Client

			client.Closer = c
			client.Reader = &Reader{c, bufio.NewReaderSize(c, 512), time.Second * 5}
			client.OnMessage = func(b []byte) (err error) {
				if _, err = client.Write([]byte("hello server !")); err != nil {
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

				if client.Writer == nil {
					block, _ := pem.Decode(pkg)

					key, err := x509.ParsePKIXPublicKey(block.Bytes)

					if err != nil {
						t.Errorf("x509 %v", err)
						return
					}

					publickey := key.(*rsa.PublicKey)
					client.Writer = NewEncoder(c, publickey)
				}

				if err := client.OnMessage(pkg); err != nil {
					t.Errorf("OnMessage %v", err)
					return
				}

				if _, err := client.Discard(len(pkg)); err != nil {
					t.Errorf("OnMessage %v", err)
					return
				}

				time.Sleep(time.Second)
			}
		})
	}

	waitGroup.Wait()
}
