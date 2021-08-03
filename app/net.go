package app

import (
	"bufio"
	"crypto/rsa"
	"net"
)

type Writer struct {
	c       net.Conn
	timeout int
	b       *Buffer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return 0, nil
}

type Reader struct {
	c       net.Conon
	timeout int
	r       *bufio.Reader
}

func (r *Reader) Read() (p []byte, err error) {
	return nil, nil
}

type Encoder struct {
	w   *Writer
	key *rsa.PublicKey
}

func (e *Encoder) Write(p []byte) (n int, err error) {
	return 0, nil
}

type Decoder struct {
	r   *Reader
	key *rsa.PrivateKey
}

func (d *Decoder) Read() (p []byte, err error) {
	return nil, nil
}

type Buffer struct {
	b []byte
}
