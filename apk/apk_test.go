package apk

import (
	"bufio"
	"bytes"
	"fmt"
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
		net.Conn
		*bufio.Reader
		*bytes.Buffer
		OnMessage func(pkg []byte) (err error)
		Write     func(pkg []byte) (n int, err error)
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

			client.Conn = c
			client.OnMessage = func(b []byte) (err error) {
				if _, err = client.Write([]byte("hello server !")); err != nil {
					return err
				}

				return nil
			}
			client.Write = func(b []byte) (n int, err error) {
				length := len(b)

				var header [2]byte

				header[0] = uint8(length >> 8)
				header[1] = uint8(length)

				if _, err := client.Buffer.Write(header[:]); err != nil {
					return 0, err
				}

				if _, err := client.Buffer.Write(b); err != nil {
					return 0, err
				}

				return client.Conn.Write(client.Buffer.Bytes())
			}

			defer func() {
				if err := client.Close(); err != nil {
					t.Errorf("close err %v", err)
				}
			}()

			if _, err := client.Write([]byte("hello server!!!")); err != nil {
				t.Errorf("Send %v", err)
				return
			}

			var header [2]byte

			for {
				if _, err := client.Reader.Read(header[:]); err != nil {
					t.Errorf("OnMessage %v", err)
					return
				}

				length := uint16(uint16(header[0])<<8 | uint16(header[1]))

				pkg, err := client.Reader.Peek(int(length))

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
