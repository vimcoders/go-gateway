package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/vimcoders/go-driver"
)

const (
	Version = 88
)

type Session struct {
	*Buffer
	net.Conn
	id               int64
	v                map[interface{}]interface{}
	PushMessageQuene chan driver.Message
	OnMessage        func(pkg driver.Message) (err error)
}

func (s *Session) WaitMessage() (err error) {
	for {
		select {
		case pkg := <-s.PushMessageQuene:
			if pkg == nil {
				return errors.New("shutdown")
			}

			if err := s.Push(pkg); err != nil {
				return err
			}
		default:
			pkg, err := s.Pull()

			if err != nil {
				return err
			}

			if err := s.OnMessage(pkg); err != nil {
				return err
			}

			time.Sleep(time.Second)
		}
	}
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

	const header = 4
	length := len(b)

	buf := s.Take(length + header)

	copy(buf[header:], b)

	buf[0] = uint8(Version >> 8)
	buf[1] = uint8(Version)
	buf[2] = uint8(length >> 8)
	buf[3] = uint8(length)

	if err := s.SetDeadline(time.Now().Add(time.Millisecond * timeout)); err != nil {
		return err
	}

	if _, err := s.Conn.Write(buf); err != nil {
		return err
	}

	return nil
}

func (s *Session) Pull() (driver.Message, error) {
	if err := s.SetDeadline(time.Now().Add(time.Millisecond * timeout)); err != nil {
		return nil, err
	}

	buffer := s.Take(DefaultBufferSize)

	if _, err := s.Read(buffer); err != nil {
		return nil, err
	}

	version := int(uint32(buffer[0])<<8 | uint32(buffer[1]))

	if version != Version {
		return nil, errors.New(fmt.Sprintf("unknow version %v", version))
	}

	length := int(uint32(buffer[2])<<8 | uint32(buffer[3]))

	const header = 4

	return NewDecoder(privateKey, buffer[header:length+header]), nil
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

func Handle(ctx context.Context, c net.Conn, pkg driver.Message) (err error) {
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
		Buffer:           NewBuffer(),
		v:                make(map[interface{}]interface{}),
	}

	s.OnMessage = func(pkg driver.Message) (err error) {
		b, err := pkg.ToBytes()

		if err != nil {
			logger.Error("OnMessage %v", err)
			return err
		}

		logger.Info("OnMessage %v..", string(b))

		if err := s.Send(NewEncoder(nil, b)); err != nil {
			return err
		}

		return nil
	}

	defer s.Close()

	if err := s.Push(pkg); err != nil {
		return err
	}

	return s.WaitMessage()
}
