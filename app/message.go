package app

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/vimcoders/go-driver"
)

type Encoder struct {
	b   []byte
	key *rsa.PublicKey
}

func (e *Encoder) ToBytes() (b []byte, err error) {
	if e.key == nil {
		return e.b, nil
	}

	return rsa.EncryptPKCS1v15(rand.Reader, e.key, e.b)
}

func NewEncoder(key *rsa.PublicKey, b []byte) driver.Message {
	return &Encoder{b, key}
}

type Decoder struct {
	b   []byte
	key *rsa.PrivateKey
}

func (d *Decoder) ToBytes() (b []byte, err error) {
	if d.key == nil {
		return d.b, nil
	}

	return rsa.DecryptPKCS1v15(rand.Reader, d.key, d.b)
}

func NewDecoder(k *rsa.PrivateKey, b []byte) driver.Message {
	return &Decoder{b, k}
}
