package app

import (
	"context"

	"github.com/vimcoders/pb"
	"google.golang.org/protobuf/proto"
)

func (s *Session) Login(ctx context.Context, b []byte) (proto.Message, error) {
	var login pb.Login

	if err := proto.Unmarshal(b, &login); err != nil {
		return nil, nil
	}

	Result, err := grpc.Login(ctx, &pb.GrpcLogin{Login: &login})

	if err != nil {
		return nil, err
	}

	return Result.Result, nil
}
