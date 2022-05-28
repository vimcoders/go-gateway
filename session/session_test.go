package session

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-gateway/pb"
	"golang.org/x/net/websocket"
	"google.golang.org/protobuf/proto"
)

func init_tcp() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	go func() {
		if e := recover(); e != nil {
			panic(e)
		}
		for {
			c, err := l.Accept()
			if err != nil {
				panic(err)
				continue
			}
			s := &Session{
				Conn: &driver.Conn{Conn: c, C: make(chan []byte, 1)},
			}
			s.OnMessage = func(b []byte) error {
				var req pb.Request
				if err := proto.Unmarshal(b, &req); err != nil {
					return err
				}
				s.Send(&pb.Response{Body: req.Body})
				return nil
			}
			go s.Pull()
			go s.Push()
		}
	}()
}

func init_websocket() {
	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		s := &Session{
			Conn: &driver.Conn{Conn: ws, C: make(chan []byte, 1)},
		}
		s.OnMessage = func(b []byte) error {
			var req pb.Request
			if err := proto.Unmarshal(b, &req); err != nil {
				return err
			}
			s.Send(&pb.Response{Body: req.Body})
			return nil
		}
		go s.Push()
		s.Pull()
	}))
	go func() {
		if e := recover(); e != nil {
			panic(e)
		}
		if err := http.ListenAndServe(":8889", nil); err != nil {
			panic(err)
		}
	}()
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	m.Run()
	fmt.Println("end")
}

// 测试发送消息
func TestTcp(t *testing.T) {
	init_tcp()
	var waitGroup sync.WaitGroup
	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		c, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			t.Error(err)
			return
		}
		s := &Session{
			Conn: &driver.Conn{Conn: c, C: make(chan []byte, 1)},
		}
		s.OnMessage = func(b []byte) error {
			t.Log(b)
			return nil
		}
		go s.Pull()
		go s.Push()
		go func() {
			defer waitGroup.Done()
			for k := 0; k < 100; k++ {
				s.Send(&pb.Request{Body: "fjafjkajf;laksjdfaskjf"})
				time.Sleep(time.Second)
			}
		}()
	}
	waitGroup.Wait()
}

// 测试发送消息
func TestWebSocket(t *testing.T) {
	init_websocket()
	var waitGroup sync.WaitGroup
	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		ws, err := websocket.Dial("ws://localhost:8889/ws", "", "http://localhost/")
		if err != nil {
			t.Error(err)
			return
		}
		s := &Session{
			Conn: &driver.Conn{Conn: ws, C: make(chan []byte, 1)},
		}
		s.OnMessage = func(b []byte) error {
			t.Log(b)
			return nil
		}
		go s.Pull()
		go s.Push()
		go func() {
			defer waitGroup.Done()
			defer s.Close()
			for k := 0; k < 100; k++ {
				s.Send(&pb.Request{Body: "fjafjkajf;laksjdfaskjf"})
				time.Sleep(time.Second)
			}
		}()
	}
	waitGroup.Wait()
}
