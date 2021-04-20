package app

import (
	"bufio"
	"context"
	"errors"
	"net"
	"time"

	"github.com/vimcoders/go-driver"
	"google.golang.org/protobuf/proto"
)

const (
	Version           = 1
	HeaderLength      = 4
	DefaultBufferSize = 128
)

type Buffer struct {
	buf []byte
}

func (b *Buffer) Take(n int) []byte {
	if n < len(b.buf) {
		return b.buf[:n]
	}

	return make([]byte, n)
}

func NewBuffer() *Buffer {
	return NewBufferSize(DefaultBufferSize)
}

func NewBufferSize(n int) *Buffer {
	return &Buffer{
		buf: make([]byte, n),
	}
}

type Encoder struct {
	proto.Message
}

func (e *Encoder) ToBytes() (b []byte, errr error) {
	return proto.Marshal(e.Message)
}

func NewEncoder(msg proto.Message) driver.Message {
	return &Encoder{msg}
}

type Decoder struct {
	b []byte
}

func (d *Decoder) ToBytes() (b []byte, err error) {
	return d.b, nil
}

func NewDecoder(b []byte) driver.Message {
	return &Decoder{b}
}

type Conn struct {
	net.Conn
	OnMessage        func(pkg driver.Message) (err error)
	OnClose          func(e interface{})
	PushMessageQuene chan driver.Message
}

func (c *Conn) Write(pkg driver.Message) (err error) {
	defer func() {
		if e := recover(); e != nil {
			c.OnClose(e)
			return
		}
	}()

	c.PushMessageQuene <- pkg

	return nil
}

func (c *Conn) Push(ctx context.Context) (err error) {
	defer func() {
		close(c.PushMessageQuene)

		if e := recover(); e != nil {
			c.OnClose(e)
			return
		}

		c.OnClose(err)
	}()

	buffer := NewBuffer()

	for {
		select {
		case <-ctx.Done():
			return errors.New("shut down")
		default:
		}

		pkg, ok := <-c.PushMessageQuene

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

		if err := c.SetWriteDeadline(time.Now().Add(time.Second * 5)); err != nil {
			return err
		}

		if _, err := c.Conn.Write(buf); err != nil {
			return err
		}
	}
}

func (c *Conn) Pull(ctx context.Context) (err error) {
	defer func() {
		if e := recover(); e != nil {
			c.OnClose(e)
			return
		}

		c.OnClose(err)
	}()

	reader := bufio.NewReaderSize(c.Conn, DefaultBufferSize)

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

		if err := c.OnMessage(NewDecoder(buf[len(header):])); err != nil {
			return err
		}

		if _, err := reader.Discard(len(buf)); err != nil {
			return err
		}
	}
}

type Session struct {
	id int64
	driver.Writer
	v map[interface{}]interface{}
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

func (s *Session) Close() error {
	return nil
}

func Handle(ctx context.Context, c net.Conn) driver.Session {
	s := Session{
		v: make(map[interface{}]interface{}),
	}

	conn := &Conn{
		Conn: c,
		OnMessage: func(pkg driver.Message) (err error) {
			return nil
		},
		OnClose: func(e interface{}) {
			if e != nil {
				logger.Error("OnClose %v", e)
			}

			if err := c.Close(); err != nil {
				logger.Error("session err %v", err.Error())
			}

			if err := s.Close(); err != nil {
				logger.Error("session err %v", err.Error())
			}
		},
		PushMessageQuene: make(chan driver.Message),
	}

	go conn.Pull(ctx)
	go conn.Push(ctx)

	s.Writer = conn

	return &s
}
