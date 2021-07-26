package app

import (
	"crypto/rand"
	"crypto/rsa"
	"net"
	"time"

	"github.com/vimcoders/go-driver"
)

type Encoder struct {
	b   []byte
	key *rsa.PublicKey
	*Buffer
	net.Conn
}

func (e *Encoder) ToBytes() (b []byte, err error) {
	if e.key == nil {
		return e.b, nil
	}

	return rsa.EncryptPKCS1v15(rand.Reader, e.key, e.b)
}

func (e *Encoder) Write(b []byte) error {
	encoder, err := rsa.EncryptPKCS1v15(rand.Reader, e.key, b)

	if err != nil {
		return err
	}

	const header = 4
	length := len(encoder)

	buf := s.Take(length + header)

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

func NewEncoder(key *rsa.PublicKey, b []byte) driver.Message {
	return &Encoder{b, key}
}

type Decoder struct {
	b   []byte
	key *rsa.PrivateKey
	net.Conn
}

func (d *Decoder) ToBytes() (b []byte, err error) {
	if d.key == nil {
		return d.b, nil
	}

	return rsa.DecryptPKCS1v15(rand.Reader, d.key, d.b)
}

func (d *Decoder) Read(b []byte) (n int, err error) {
	return nil, nil
}

func NewDecoder(k *rsa.PrivateKey, b []byte) driver.Message {
	return &Decoder{b, k}
}
