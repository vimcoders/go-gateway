package apk

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net"
)

type Session struct {
	io.Closer
	io.Writer
	*Reader
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
		Reader: NewReader(c),
		v:      make(map[interface{}]interface{}),
	}

	defer s.Close()

	if _, err := s.Write(pem.EncodeToMemory(block)); err != nil {
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

		if _, err := s.Discard(len(p)); err != nil {
			return err
		}
	}
}
