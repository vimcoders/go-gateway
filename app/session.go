package app

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net"
	"time"

	"github.com/vimcoders/go-driver"
)

const (
	Version = 1
)

type Session struct {
	*Buffer
	net.Conn
	id               int64
	v                map[interface{}]interface{}
	PushMessageQuene chan driver.Message
	key              *rsa.PrivateKey
	OnMessage        func(pkg driver.Message) (err error)
}

func (s *Session) Send(pkg driver.Message) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("write recoder %v", e)
		}

		if err != nil {
			logger.Error("write err %v", err)
		}
	}()

	s.PushMessageQuene <- pkg

	return nil
}

func (s *Session) Push(pkg driver.Message) (err error) {
	b, err := pkg.ToBytes()

	if err != nil {
		return err
	}

	length := len(b)

	const header = 4

	buffer := make([]byte, header)

	buffer[0] = Version
	buffer[1] = uint8(length >> 16)
	buffer[2] = uint8(length >> 8)
	buffer[3] = uint8(length)

	buf := s.Take(length + header)

	copy(buf, buffer)
	copy(buf[header:], b)

	if err := s.SetDeadline(time.Now().Add(time.Millisecond * timeout)); err != nil {
		return err
	}

	if _, err := s.Conn.Write(buf); err != nil {
		return err
	}

	return nil
}

func (s *Session) Pull() (driver.Message, error) {
	buffer := s.Buffer.Buffer()

	if err := s.SetDeadline(time.Now().Add(time.Millisecond * timeout)); err != nil {
		return nil, err
	}

	if _, err := s.Read(buffer); err != nil {
		return nil, err
	}

	if buffer[0] != Version {
		return nil, errors.New("unknow version")
	}

	length := int(uint32(uint32(buffer[1])<<16 | uint32(buffer[2])<<8 | uint32(buffer[3])))
	const header = 4

	return NewDecoder(s.key, buffer[header:length+header]), nil
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

func (s *Session) Close() (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("session close %v", e)
		}

		if err != nil {
			logger.Error("session close %v", err)
		}
	}()

	close(s.PushMessageQuene)
	return s.Conn.Close()
}

func (s *Session) Handshake() (err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)

	if err != nil {
		return err
	}

	b, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}

	decoder := NewEncoder(nil, pem.EncodeToMemory(block))

	if err := s.Send(decoder); err != nil {
		return err
	}

	s.key = privateKey

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
		Conn:             c,
		PushMessageQuene: make(chan driver.Message, 1),
		OnMessage: func(pkg driver.Message) (err error) {
			b, err := pkg.ToBytes()

			if err != nil {
				logger.Error("OnMessage %v", err)
				return err
			}

			logger.Info("OnMessage %v..", string(b))

			return nil
		},
		v:      make(map[interface{}]interface{}),
		Buffer: NewBuffer(),
	}

	defer s.Close()

	if err := s.Handshake(); err != nil {
		return err
	}

	for {
		select {
		case pkg := <-s.PushMessageQuene:
			if err := s.Push(pkg); err != nil {
				return err
			}
		default:
			pkg, err := s.Pull()

			if err != nil {
				return err
			}

			s.OnMessage(pkg)
		}
	}
}
