package handler

import (
	"awesomeProject/src/message"
	"awesomeProject/src/model"
)

type Handler interface {
	// 业务处理
	Handle(spiderInfo *model.SpiderInfo, message message.IMessage) (err error)
	// 判断是否为新数据，如不是则不会处理下载或切词等操作
	IsNewMessage(message message.IMessage) (bool, error)
}
