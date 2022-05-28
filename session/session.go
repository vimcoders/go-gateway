package session

import (
	"errors"
	"fmt"

	"github.com/vimcoders/go-driver"
	"google.golang.org/protobuf/proto"
)

type Session struct {
	Id int64
	*driver.Conn
	d map[interface{}]interface{}
}

func (s *Session) Set(key, value interface{}) error {
	s.d[key] = value
	return nil
}

func (s *Session) Get(key interface{}) interface{} {
	if v, ok := s.d[key]; ok {
		return v
	}
	return nil
}

func (s *Session) Delete(key interface{}) error {
	delete(s.d, key)
	return nil
}

func (s *Session) Push() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()
	defer s.Close()
	return s.Conn.Push()
}

func (s *Session) Pull() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()
	defer s.Close()
	return s.Conn.Pull()
}

func (s *Session) Send(m proto.Message) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	s.C <- b
	return nil
}
