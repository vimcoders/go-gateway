package lib

import (
	"net/http"
	"time"
)

var (
	addr    = ":8888"
	unix    = time.Now().Unix()
	timeout = time.Second * 15
)

func init() {
	go http.ListenAndServe(":8080", nil)
}

func Timeout() time.Time {
	return time.Now().Add(timeout)
}

func Addr() string {
	return addr
}

func Seconds() int64 {
	return time.Now().Unix() - unix
}
