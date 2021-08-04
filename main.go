package main

import (
	"runtime"

	"github.com/vimcoders/go-gateway/app"
)

func main() {
	runtime.GOMAXPROCS(3)
	app.Run()
}
