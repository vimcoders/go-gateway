package app

import (
	"context"
	"net"

	"github.com/vimcoders/go-driver"
	"github.com/vimcoders/go-lib"
)

type Session struct {
	id int64
	driver.Conn
	v map[interface{}]interface{}
}

func (s *Session) SessionID() int64 {
	return s.id
}

func (s *Session) Set(key, value interface{}) error {
	s.v[key] = value
	return nil
}

func (s *Session) Get(key interface{}) interface{} {
	return s.v[key]
}

func (s *Session) Delete(key interface{}) error {
	delete(s.v, key)
	return nil
}

func NewSession(ctx context.Context, c net.Conn) driver.Session {
	s := Session{
		v: make(map[interface{}]interface{}),
	}

	conn := &lib.Conn{
		Conn: c,
		OnMessage: func(pkg driver.Message) (err error) {
			return nil
		},
		OnClose: func(e interface{}) {
		},
		PushMessageQuene: make(chan driver.Message),
	}

	go conn.Pull(ctx)
	go conn.Push(ctx)

	s.Conn = conn

	return &s
}
