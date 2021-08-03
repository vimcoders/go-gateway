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
	"github.com/vimcoders/go-lib"
)

type Session struct {
	io.Closer
	io.Writer
	driver.Reader

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
		Reader: lib.NewReader(c, lib.NewBuffer(), time.Second*5),
		Writer: lib.NewWriter(c, lib.NewBuffer(), time.Second*5),
		v:      make(map[interface{}]interface{}),
	}

	s.OnMessage = func(p driver.Message) (err error) {
		b, err := p.ToBytes()

		if err != nil {
			logger.Error("OnMessage %v", err)
			return err
		}

		logger.Info("OnMessage %v..", string(b))

		if err := s.Writer.Write(lib.NewMessage([]byte("hello client !"))); err != nil {
			return err
		}

		return nil
	}

	defer s.Close()

	pkg := lib.NewMessage(pem.EncodeToMemory(block))

	if err := s.Writer.Write(pkg); err != nil {
		return err
	}

	for {
		p, err := s.Reader.Read()

		if err != nil {
			return err
		}

		b, err := p.ToBytes()

		if err != nil {
			return err
		}

		decoder := lib.NewDecoder(b, privateKey)

		if err := s.OnMessage(decoder); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}
}
