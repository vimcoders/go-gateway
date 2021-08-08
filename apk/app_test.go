package apk

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
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

	for i := 0; i < 1; i++ {
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
			client.Reader = NewReader(c)
			client.Writer = NewWriter(c)
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

			client.Writer.Write([]byte("hello server!!!"))

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

				if _, err := client.Discard(len(pkg)); err != nil {
					t.Errorf("OnMessage %v", err)
					return
				}
			}
		})
	}

	waitGroup.Wait()
}
