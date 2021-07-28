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
	driver.Reader
	driver.Writer

	OnMessage func(pkg driver.Message) (err error)
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
		Reader: NewDecoder(c, NewBuffer(), privateKey),
		Writer: NewWriter(c, NewBuffer()),
		v:      make(map[interface{}]interface{}),
	}

	s.OnMessage = func(p driver.Message) (err error) {
		b, err := p.ToBytes()

		if err != nil {
			logger.Error("OnMessage %v", err)
			return err
		}

		logger.Info("OnMessage %v..", string(b))

		if err := s.Writer.Write(NewMessage([]byte("hello client !"))); err != nil {
			return err
		}

		return nil
	}

	defer s.Close()

	pkg := NewMessage(pem.EncodeToMemory(block))

	if err := s.Writer.Write(pkg); err != nil {
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
