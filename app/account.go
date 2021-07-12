package app

import (
	"github.com/vimcoders/go-gateway/pb"
	"google.golang.org/protobuf/proto"
)

func (s *Session) Login(b []byte) (err error) {
	var login pb.Login

	if err := proto.Unmarshal(b, &login); err != nil {
		return err
	}

	return nil
}

func (s *Session) Register(b []byte) (err error) {
	return nil
}
