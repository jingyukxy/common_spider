package executor

import (
	"awesomeProject/src/message"
	"awesomeProject/src/model"
)


type SegmentExecutor struct {
}

func (segExec *SegmentExecutor) Execute(info *model.SpiderInfo, msg message.IMessage) message.IMessage {
	return msg
}
