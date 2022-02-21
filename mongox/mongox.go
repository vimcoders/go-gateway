package mongox

import (
	"github.com/vimcoders/go-gateway/logx"
)

//func init() {
//	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
//
//	// 连接到MongoDB
//	mgoCli, err := mongo.Connect(context.TODO(), clientOptions)
//	if err != nil {
//		log.Error("err %v", err)
//		return
//	}
//
//	// 检查连接
//	if err := mgoCli.Ping(context.TODO(), nil); err != nil {
//		log.Error("err %v", err)
//		return
//	}
//
//	Cli = mgoCli
//}

func init() {
	logx.Info("init mongox ...")
}
