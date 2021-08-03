package app

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net"
	"time"

	"github.com/vimcoders/go-driver"
)

type Writer struct {
	c       net.Conn
	b       *Buffer
	timeout time.Duration
}

func (w *Writer) Write(p []byte) (n int, err error) {
	length := len(p)

	buf := w.b.Take(length + 2)
	buf[0] = uint8(length >> 8)
	buf[1] = uint8(length)
	copy(buf[2:], p)

	if err := w.c.SetWriteDeadline(time.Now().Add(w.timeout)); err != nil {
		return 0, err
	}

	return w.c.Write(buf)
}

func NewWriter(c net.Conn) io.Writer {
	return &Writer{c, NewBuffer(), time.Second * 15}
}

type Reader struct {
	c       net.Conn
	r       *bufio.Reader
	timeout time.Duration
}

func (r *Reader) Read() (p []byte, err error) {
	header := make([]byte, 2)

	if err := r.c.SetReadDeadline(time.Now().Add(r.timeout)); err != nil {
		return nil, err
	}

	if _, err := r.r.Read(header); err != nil {
		return nil, err
	}

	length := uint16(uint16(header[0])<<8 | uint16(header[1]))

	if err := r.c.SetReadDeadline(time.Now().Add(r.timeout)); err != nil {
		return nil, err
	}

	return r.r.Peek(int(length))
}

func NewReader(c net.Conn) driver.Reader {
	return &Reader{c, bufio.NewReaderSize(c, 512), time.Second * 5}
}

type Encoder struct {
	w   *Writer
	key *rsa.PublicKey
}

func (e *Encoder) Write(p []byte) (n int, err error) {
	b, err := rsa.EncryptPKCS1v15(rand.Reader, e.key, p)

	if err != nil {
		return 0, err
	}

	return e.w.Write(b)
}

func NewEncoder(c net.Conn, k *rsa.PublicKey) io.Writer {
	return &Encoder{
		w:   &Writer{c, NewBuffer(), time.Second * 5},
		key: k,
	}
}

type Decoder struct {
	r   *Reader
	key *rsa.PrivateKey
}

func (d *Decoder) Read() (p []byte, err error) {
	b, err := d.r.Read()

	if err != nil {
		return nil, err
	}

	buf, err := rsa.DecryptPKCS1v15(rand.Reader, d.key, b)

	if err != nil {
		return
	}

	if _, err := d.r.r.Discard(len(b)); err != nil {
		return nil, err
	}

	return buf, nil
}

func NewDecoder(c net.Conn, k *rsa.PrivateKey) *Decoder {
	return &Decoder{
		r:   &Reader{c, bufio.NewReaderSize(c, 512), time.Second * 5},
		key: k,
	}
}

type Buffer struct {
	b []byte
}

func (b *Buffer) Take(n int) (p []byte) {
	if n > len(b.b) {
		return make([]byte, n)
	}

	return b.b[:n]
}

func NewBuffer() *Buffer {
	return &Buffer{
		b: make([]byte, 512),
	}
}
