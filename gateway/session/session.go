package session

import (
	"bufio"
	"bytes"
	"context"
	"net"

	"github.com/vimcoders/go-gateway/gateway/log"
)

type Session struct {
	net.Conn
	*bytes.Buffer
	*bufio.Reader
}

func (s *Session) Set(key, value interface{}) error {
	return nil
}

func (s *Session) Get(key interface{}) interface{} {
	return nil
}

func (s *Session) Delete(key interface{}) error {
	return nil
}

func (s *Session) OnMessage(p []byte) error {
	log.Info("OnMessage %v..", string(p))

	if _, err := s.Write([]byte("hello client !")); err != nil {
		return err
	}

	return nil
}

func (s *Session) Write(b []byte) (n int, err error) {
	defer s.Buffer.Reset()

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
			log.Error("Handle recover %v", e)
		}

		if err != nil {
			log.Error("Handle %v", err)
		}
	}()

	s := Session{
		Conn:   c,
		Reader: bufio.NewReader(c),
		Buffer: bytes.NewBuffer(make([]byte, 512)),
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
