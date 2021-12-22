package opt

import (
	"context"
	"time"
)

var (
	Addr                = ":8888"
	Network             = "tcp"
	CloseCtx, CloseFunc = context.WithCancel(context.Background())
	Timeout             = time.Second * 15
)
