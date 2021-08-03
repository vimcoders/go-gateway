package app

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net"
	"time"

	"github.com/vimcoders/go-driver"
)

type Session struct {
	io.Closer
	io.Writer
	driver.Reader

	OnMessage func(pkg []byte) (err error)
	v         map[interface{}]interface{}
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

func Handle(ctx context.Context, c net.Conn, k *rsa.PrivateKey) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("Handle recover %v", e)
		}

		if err != nil {
			logger.Error("Handle %v", err)
		}
	}()

	b, err := x509.MarshalPKIXPublicKey(&k.PublicKey)

	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}

	s := Session{
		Closer: c,
		Writer: NewWriter(c),
		Reader: NewDecoder(c, k),
		v:      make(map[interface{}]interface{}),
	}

	s.OnMessage = func(p []byte) (err error) {
		logger.Info("OnMessage %v..", string(p))

		if _, err := s.Writer.Write([]byte("hello client !")); err != nil {
			return err
		}

		s.Close()

		return nil
	}

	defer s.Close()

	pkg := pem.EncodeToMemory(block)

	if _, err := s.Write(pkg); err != nil {
		return err
	}

	for {
		p, err := s.Reader.Read()

		if err != nil {
			return err
		}

		if err := s.OnMessage(p); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}
}
