package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	for i := 0; i < 100000; i++ {
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 64)

			if err != nil {
				t.Error(err)
				return
			}

			b, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

			if err != nil {
				t.Error(err)
				return
			}

			block := &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: b,
			}

			t.Log(string(pem.EncodeToMemory(block)))
		})
	}
}
