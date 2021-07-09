package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestHandle(t *testing.T) {
	tests := []struct {
		name string
		bits int
	}{
		{"case 1", 66},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := rsa.GenerateKey(rand.Reader, tt.bits)

			if err != nil {
				t.Error(err)
				return
			}

			bytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)

			if err != nil {
				t.Error(err)
				return
			}

			block := pem.Block{Type: "RSA Public Key", Bytes: bytes}

			t.Log(string(pem.EncodeToMemory(&block)))
		})
	}
}
