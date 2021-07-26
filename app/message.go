package app

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/vimcoders/go-driver"
)

type Message struct {
	message []byte
}

func (m *Message) ToBytes() (b []byte, err error) {
	return m.message, nil
}

func NewMessage(b []byte) driver.Message {
	return &Message{b}
}

type Encoder struct {
	net.Conn
	buffer *Buffer
	key    *rsa.PublicKey
}

func (e *Encoder) Write(pkg driver.Message) (err error) {
	b, err := pkg.ToBytes()

	if err != nil {
		return err
	}

	encoder, err := rsa.EncryptPKCS1v15(rand.Reader, e.key, b)

	if err != nil {
		return err
	}

	const header = 4
	length := len(encoder)

	buf := e.buffer.Take(length + header)

	copy(buf[header:], encoder)

	buf[0] = uint8(Version >> 8)
	buf[1] = uint8(Version)
	buf[2] = uint8(length >> 8)
	buf[3] = uint8(length)

	if err := e.SetDeadline(time.Now().Add(time.Millisecond * timeout)); err != nil {
		return err
	}

	if _, err := e.Conn.Write(buf); err != nil {
		return err
	}

	return nil
}

func NewEncoder(c net.Conn, b *Buffer, k *rsa.PublicKey) *Encoder {
	return &Encoder{c, b, k}
}

type Writer struct {
	net.Conn
	b *Buffer
}

func (w *Writer) Write(pkg driver.Message) (err error) {
	b, err := pkg.ToBytes()

	if err != nil {
		return err
	}

	const header = 4
	length := len(b)

	buf := w.b.Take(length + header)

	copy(buf[header:], b)

	buf[0] = uint8(Version >> 8)
	buf[1] = uint8(Version)
	buf[2] = uint8(length >> 8)
	buf[3] = uint8(length)

	if err := w.SetDeadline(time.Now().Add(time.Millisecond * timeout)); err != nil {
		return err
	}

	if _, err := w.Conn.Write(buf); err != nil {
		return err
	}

	return nil
}

func NewWriter(c net.Conn, b *Buffer) *Writer {
	return &Writer{c, b}
}

type Decoder struct {
	net.Conn
	buffer *Buffer
	key    *rsa.PrivateKey
}

func (d *Decoder) Read() (pkg driver.Message, err error) {
	if err := d.SetDeadline(time.Now().Add(time.Millisecond * timeout)); err != nil {
		return nil, err
	}

	buffer := d.buffer.Take(DefaultBufferSize)

	if _, err := d.Conn.Read(buffer); err != nil {
		return nil, err
	}

	version := int(uint32(buffer[0])<<8 | uint32(buffer[1]))

	if version != Version {
		return nil, errors.New(fmt.Sprintf("unknow version %v", version))
	}

	length := int(uint32(buffer[2])<<8 | uint32(buffer[3]))

	const header = 4

	body := buffer[header : header+length]

	decoder, err := rsa.DecryptPKCS1v15(rand.Reader, d.key, body)

	if err != nil {
		return nil, err
	}

	return &Message{decoder}, nil
}

func NewDecoder(c net.Conn, b *Buffer, k *rsa.PrivateKey) *Decoder {
	return &Decoder{c, b, k}
}

type Reader struct {
	net.Conn
	buffer *Buffer
}

func (r *Reader) Read() (pkg driver.Message, err error) {
	if err := r.SetDeadline(time.Now().Add(time.Millisecond * timeout)); err != nil {
		return nil, err
	}

	buffer := r.buffer.Take(DefaultBufferSize)

	if _, err := r.Conn.Read(buffer); err != nil {
		return nil, err
	}

	version := int(uint32(buffer[0])<<8 | uint32(buffer[1]))

	if version != Version {
		return nil, errors.New(fmt.Sprintf("unknow version %v", version))
	}

	length := int(uint32(buffer[2])<<8 | uint32(buffer[3]))

	const header = 4

	body := buffer[header : header+length]

	return &Message{body}, nil
}

func NewReader(c net.Conn, b *Buffer) *Reader {
	return &Reader{c, b}
}
