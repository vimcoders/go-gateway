package bot

import (
	"io"
	"net"
	"sync"

	"github.com/vimcoders/go-driver"
)

type Bot struct {
	io.Closer
	io.Writer
	*driver.Reader
	OnMessage func(pkg []byte) (err error)
}

func (b *Bot) Login() error {
	return nil
}

func NewBot(c net.Conn) *Bot {
	var bot Bot

	bot.Closer = c
	bot.Writer = driver.NewWriter(c)
	bot.Reader = driver.NewReader(c)

	return &bot
}

func Login() {
	var waitGroup sync.WaitGroup

	for i := 0; i < 1; i++ {
		waitGroup.Add(1)

		go func() (err error) {
			defer waitGroup.Done()

			c, err := net.Dial("tcp", ":8888")

			if err != nil {
				return
			}

			var bot Bot

			bot.Closer = c
			bot.Reader = driver.NewReader(c)
			bot.Writer = driver.NewWriter(c)

			bot.Login()

			for {
				pkg, err := bot.Read()

				if err != nil {
					return err
				}

				if err := bot.OnMessage(pkg); err != nil {
					return err
				}

				if _, err := bot.Discard(len(pkg)); err != nil {
					return err
				}
			}
		}()
	}

	waitGroup.Wait()
}
