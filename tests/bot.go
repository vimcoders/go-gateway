package tests

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net"
)

type Bot struct {
	net.Conn
	*bytes.Buffer
	*bufio.Reader
}

func (b *Bot) Write(pkg []byte) (n int, err error) {
	defer b.Buffer.Reset()

	length := len(pkg)

	var header [2]byte

	header[0] = uint8(length >> 8)
	header[1] = uint8(length)

	if _, err := b.Buffer.Write(header[:]); err != nil {
		return 0, err
	}

	if _, err := b.Buffer.Write(pkg); err != nil {
		return 0, err
	}

	fmt.Println(b.Buffer.Bytes())

	return b.Conn.Write(b.Buffer.Bytes())
}

func (b *Bot) Read() (pkg []byte, err error) {
	header := make([]byte, 2)

	for {
		if _, err := b.Reader.Read(header); err != nil {
			return nil, err
		}

		length := uint16(uint16(header[0])<<8 | uint16(header[1]))

		pkg, err := b.Reader.Peek(int(length))

		if err != nil {
			return nil, err
		}

		if err := b.OnMessage(pkg); err != nil {
			return nil, err
		}

		if _, err := b.Discard(len(pkg)); err != nil {
			return nil, err
		}

		return pkg, nil
	}
}

func (b *Bot) OnMessage(pkg []byte) (err error) {
	return nil
}

func NewBot(ctx context.Context, c net.Conn) *Bot {
	return &Bot{
		Conn:   c,
		Reader: bufio.NewReader(c),
		Buffer: bytes.NewBuffer(make([]byte, 512)),
	}
}

func (b *Bot) Login() (err error) {
	if _, err := b.Write([]byte("login")); err != nil {
		return err
	}

	if _, err := b.Read(); err != nil {
		return err
	}

	fmt.Println("login")

	return nil
}

func (b *Bot) Register() (err error) {
	if _, err := b.Write([]byte("register")); err != nil {
		return err
	}

	if _, err := b.Read(); err != nil {
		return err
	}

	fmt.Println("register")

	return nil
}
