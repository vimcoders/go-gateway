package bot

import (
	"fmt"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	m.Run()
	fmt.Println("end")
}

func TestLogin(t *testing.T) {
	var waitGroup sync.WaitGroup

	for i := 0; i < 1000; i++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			bot := NewBot()

			if err := bot.Login(); err != nil {
				return
			}

			for {
				pkg, err := bot.Read()

				if err != nil {
					continue
				}

				if err := bot.Login(); err != nil {
					continue
				}

				if _, err := bot.Discard(len(pkg)); err != nil {
					continue
				}
			}
		}()
	}
}

func TestRegister(t *testing.T) {
	var waitGroup sync.WaitGroup

	for i := 0; i < 1000; i++ {
		go func() {
			defer waitGroup.Done()

			bot := NewBot()

			if err := bot.Register(); err != nil {
				return
			}

			for {
				pkg, err := bot.Read()

				if err != nil {
					continue
				}

				if err := bot.Register(); err != nil {
					continue
				}

				if _, err := bot.Discard(len(pkg)); err != nil {
					continue
				}
			}
		}()
	}
}
