package apk

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net"
	"time"
)

type Writer struct {
	net.Conn
	*Buffer
	t time.Duration
}

func (w *Writer) Write(p []byte) (n int, err error) {
	length := len(p)

	buf := w.Take(length + 2)
	buf[0] = uint8(length >> 8)
	buf[1] = uint8(length)
	copy(buf[2:], p)

	if err := w.SetWriteDeadline(time.Now().Add(w.t)); err != nil {
		return 0, err
	}

	return w.Conn.Write(buf)
}

func NewWriter(c net.Conn) io.Writer {
	return &Writer{c, NewBuffer(), time.Second * 15}
}

type Reader struct {
	net.Conn
	*bufio.Reader
	t time.Duration
}

func (r *Reader) Read() (p []byte, err error) {
	header := make([]byte, 2)

	if err := r.SetReadDeadline(time.Now().Add(r.t)); err != nil {
		return nil, err
	}

	if _, err := r.Reader.Read(header); err != nil {
		return nil, err
	}

	length := uint16(uint16(header[0])<<8 | uint16(header[1]))

	if err := r.SetReadDeadline(time.Now().Add(r.t)); err != nil {
		return nil, err
	}

	return r.Reader.Peek(int(length))
}

func NewReader(c net.Conn) *Reader {
	return &Reader{c, bufio.NewReaderSize(c, 512), time.Second * 5}
}

type Encoder struct {
	*Writer
	key *rsa.PublicKey
}

func (e *Encoder) Write(p []byte) (n int, err error) {
	b, err := rsa.EncryptPKCS1v15(rand.Reader, e.key, p)

	if err != nil {
		return 0, err
	}

	return e.Writer.Write(b)
}

func NewEncoder(c net.Conn, k *rsa.PublicKey) io.Writer {
	return &Encoder{
		Writer: &Writer{c, NewBuffer(), time.Second * 5},
		key:    k,
	}
}

type Decoder struct {
	*Reader
	key *rsa.PrivateKey
}

func (d *Decoder) Read() (p []byte, err error) {
	b, err := d.Reader.Read()

	if err != nil {
		return nil, err
	}

	buf, err := rsa.DecryptPKCS1v15(rand.Reader, d.key, b)

	if err != nil {
		return
	}

	if _, err := d.Discard(len(b)); err != nil {
		return nil, err
	}

	return buf, nil
}

func NewDecoder(c net.Conn, k *rsa.PrivateKey) *Decoder {
	return &Decoder{
		Reader: &Reader{c, bufio.NewReaderSize(c, 512), time.Second * 5},
		key:    k,
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
