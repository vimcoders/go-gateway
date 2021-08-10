package bot

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"time"

	"github.com/vimcoders/go-driver"
)

type Bot struct {
	io.Closer
	io.Writer
	driver.Reader
	OnMessage func(pkg []byte) (err error)
}

func (b *Bot) Login() error {
	return nil
}

func NewBot(c net.Conn) *Bot {
	var bot Bot

	bot.Closer = c
	bot.Reader = driver.NewReader(c, bufio.NewReaderSize(c, 1024), time.Second*15)
	bot.Writer = driver.NewWriter(c, bytes.NewBuffer(make([]byte, 256)), time.Second*15)

	return &bot
}
