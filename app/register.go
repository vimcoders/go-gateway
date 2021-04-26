package app

import (
	"context"

	"github.com/vimcoders/pb"
	"google.golang.org/protobuf/proto"
)

func (s *Session) Register(ctx context.Context, b []byte) (proto.Message, error) {
	var register pb.Register

	if err := proto.Unmarshal(b, &register); err != nil {
		return nil, nil
	}

	Result, err := grpc.Register(ctx, &pb.GrpcRegister{Register: &register})

	if err != nil {
		return nil, err
	}

	return Result.Result, nil

}
