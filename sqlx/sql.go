package sqlx

import (
	"sync"

	"github.com/vimcoders/go-driver"
)

var connector driver.Connector

//func init() {
//	c, err := sqlx.Connect(nil)
//
//	if err != nil {
//		log.Error("err %v", err)
//	}
//
//	connector = c
//}

func Init(wg *sync.WaitGroup) {
	defer wg.Done()
}
