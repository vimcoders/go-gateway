package tests

import (
	"context"
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
	var waitGroup sync.WaitGroup

	for i := 0; i < 1; i++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			c, err := net.Dial("tcp", "127.0.0.1:8888")

			if err != nil {
				t.Log(err)
			}

			defer c.Close()

			bot := NewBot(context.Background(), c)

			for {
				if err := bot.Login(); err != nil {
					continue
				}
			}
		}()
	}

	waitGroup.Wait()
}

func TestRegister(t *testing.T) {
	var waitGroup sync.WaitGroup

	for i := 0; i < 1000; i++ {
		go func() {
			defer waitGroup.Done()

			c, err := net.Dial("tcp", "127.0.0.1:8888")

			if err != nil {
				t.Log(err)
			}

			bot := NewBot(context.Background(), c)

			for {
				if err := bot.Register(); err != nil {
					continue
				}
			}
		}()
	}

	waitGroup.Wait()
}
