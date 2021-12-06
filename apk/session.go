package apk

import (
	"bufio"
	"bytes"
	"context"
	"net"
)

type Session struct {
	net.Conn
	*bytes.Buffer
	*bufio.Reader
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

	if _, err := s.Write([]byte("hello client !")); err != nil {
		return err
	}

	return nil
}

func (s *Session) Write(b []byte) (n int, err error) {
	length := len(b)

	var header [2]byte

	header[0] = uint8(length >> 8)
	header[1] = uint8(length)

	if _, err := s.Buffer.Write(header[:]); err != nil {
		return 0, err
	}

	if _, err := s.Buffer.Write(b); err != nil {
		return 0, err
	}

	return s.Conn.Write(s.Buffer.Bytes())
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
		Conn:   c,
		Reader: bufio.NewReader(c),
		Buffer: bytes.NewBuffer(make([]byte, 512)),
		v:      make(map[interface{}]interface{}),
	}

	defer s.Close()

	header := make([]byte, 2)

	for {
		if _, err := s.Reader.Read(header); err != nil {
			return err
		}

		length := uint16(uint16(header[0])<<8 | uint16(header[1]))

		b, err := s.Reader.Peek(int(length))

		if err != nil {
			return err
		}

		if err := s.OnMessage(b); err != nil {
			return err
		}

		if _, err := s.Discard(len(b)); err != nil {
			return err
		}
	}
}
