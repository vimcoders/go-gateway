package app

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net"
	"runtime/debug"
	"time"

	"github.com/vimcoders/go-driver"
)

const (
	Version      = 1
	HeaderLength = 4
)

type Session struct {
	net.Conn
	id               int64
	v                map[interface{}]interface{}
	PushMessageQuene chan driver.Message
	key              *rsa.PrivateKey
	OnMessage        func(pkg driver.Message) (err error)
}

func (s *Session) Write(pkg driver.Message) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("write recoder %v debug %v", e, string(debug.Stack()))
		}

		if err != nil {
			logger.Error("write err %v", err)
		}
	}()

	s.PushMessageQuene <- pkg

	return nil
}

func (s *Session) Push(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("push recoder %v debug %v", e, string(debug.Stack()))
		}

		if err != nil {
			logger.Error("push err %v", err)
		}

		s.Close()
	}()

	buffer := NewBuffer()

	for {
		select {
		case <-ctx.Done():
			return errors.New("shut down")
		default:
		}

		pkg, ok := <-s.PushMessageQuene

		if !ok {
			return errors.New("shutdown")
		}

		b, err := pkg.ToBytes()

		if err != nil {
			return err
		}

		header := make([]byte, HeaderLength)
		length := len(b)

		header[0] = Version
		header[1] = uint8(length >> 16)
		header[2] = uint8(length >> 8)
		header[3] = uint8(length)

		buf := buffer.Take(len(header) + len(b))
		copy(buf, header)
		copy(buf[len(header):], b)

		if err := s.SetWriteDeadline(time.Now().Add(time.Second * 5)); err != nil {
			return err
		}

		if _, err := s.Conn.Write(buf); err != nil {
			return err
		}
	}
}

func (s *Session) Pull(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("pull recoder %v debug %v", e, string(debug.Stack()))
		}

		if err != nil {
			logger.Error("pull err %v", err)
		}
	}()

	reader := bufio.NewReaderSize(s.Conn, DefaultBufferSize)

	for {
		select {
		case <-ctx.Done():
			return errors.New("shut down")
		default:
		}

		header, err := reader.Peek(HeaderLength)

		if err != nil {
			return err
		}

		if header[0] != Version {
			return errors.New("Version is Unknown")
		}

		length := int(uint32(uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3])))

		buf, err := reader.Peek(length + len(header))

		if err != nil {
			return err
		}

		decoder := NewDecoder(s.key, buf[len(header):])

		if err := s.OnMessage(decoder); err != nil {
			return err
		}

		if _, err := reader.Discard(len(buf)); err != nil {
			return err
		}
	}
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
			logger.Error("session close %v debug %v", e, string(debug.Stack()))
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

	if err := s.Write(decoder); err != nil {
		return err
	}

	s.key = privateKey

	return nil
}

func Handle(ctx context.Context, c net.Conn) driver.Session {
	s := Session{
		Conn:             c,
		v:                make(map[interface{}]interface{}),
		PushMessageQuene: make(chan driver.Message),
		OnMessage: func(pkg driver.Message) (err error) {
			b, err := pkg.ToBytes()

			if err != nil {
				logger.Error("OnMessage %v", err)
				return err
			}

			logger.Info("OnMessage %v..", string(b))

			return c.Close()
		},
	}

	go s.Pull(ctx)
	go s.Push(ctx)

	if err := s.Handshake(); err != nil {
		logger.Error("Handshake %v", err)
		s.Close()
		return nil
	}

	return &s
}
