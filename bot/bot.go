package bot

import (
	"bufio"
	"bytes"
	"net"
)

type Bot struct {
	c net.Conn
	*bufio.Reader
	*bytes.Buffer
}

func (b *Bot) Login() error {
	if _, err := b.c.Write([]byte("login")); err != nil {
		return err
	}

	return nil
}

func (b *Bot) Register() error {
	if _, err := b.Write([]byte("login")); err != nil {
		return err
	}

	return nil
}

func NewBot() *Bot {
	c, err := net.Dial("tcp", ":8888")

	if err != nil {
		return nil
	}

	return &Bot{c: c}
}

func (b *Bot) Read() (pkg []byte, err error) {
	return nil, nil
}

func (b *Bot) Write(pkg []byte) (n int, err error) {
	return 0, nil
}
