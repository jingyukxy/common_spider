package executor

import (
	"awesomeProject/src/message"
	"awesomeProject/src/model"
)

type Manager struct {
	filters []Executor
}

func (manager *Manager) RegisterFilter(executor Executor) {
	manager.filters = append(manager.filters, executor)
}

func (manager *Manager) Execute(info *model.SpiderInfo, msg message.IMessage) {
	var newMsg message.IMessage
	for i, filter := range manager.filters {
		if i == 0 {
			newMsg = filter.Execute(info, msg)
		} else {
			filter.Execute(info, newMsg)
		}
	}
}

type Executor interface {
	Execute(*model.SpiderInfo, message.IMessage) message.IMessage
}

var _ Executor = (*SegmentExecutor)(nil)
var _ Executor = (*TranslateExecutor)(nil)
var _ Executor = (*DownloadExecutor)(nil)
