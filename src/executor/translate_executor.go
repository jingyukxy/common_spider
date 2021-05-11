package executor

import (
	"awesomeProject/src/message"
	"awesomeProject/src/model"
)

type TranslateExecutor struct {
}

func (trans *TranslateExecutor) Execute(info *model.SpiderInfo, message message.IMessage) message.IMessage {
	return message
}
