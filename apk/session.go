package apk

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"time"

	"github.com/vimcoders/go-driver"
)

type Session struct {
	io.Closer
	io.Writer
	driver.Reader
	v map[interface{}]interface{}
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

func (s *Session) OnMessage(p []byte) error {
	logger.Info("OnMessage %v..", string(p))

	if _, err := s.Writer.Write([]byte("hello client !")); err != nil {
		return err
	}

	return nil
}

func Handle(ctx context.Context, c net.Conn) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("Handle recover %v", e)
		}

		if err != nil {
			logger.Error("Handle %v", err)
		}
	}()

	s := Session{
		Closer: c,
		Writer: driver.NewWriter(c, bytes.NewBuffer(make([]byte, 1024)), time.Second*15),
		Reader: driver.NewReader(c, bufio.NewReaderSize(c, 256), time.Second*15),
		v:      make(map[interface{}]interface{}),
	}

	defer s.Close()

	for {
		p, err := s.Read()

		if err != nil {
			return err
		}

		if err := s.OnMessage(p); err != nil {
			return err
		}

		if _, err := s.Discard(len(p)); err != nil {
			return err
		}
	}
}
