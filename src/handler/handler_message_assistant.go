package handler

import (
	"awesomeProject/src/message"
	"sync"
)

var instance *assistant
var lock sync.Once

type assistant struct {
	processorMap map[int]*msgCodeProcessor
}

func GetAssistantInstance() *assistant {
	lock.Do(func() {
		instance = &assistant{
			processorMap: make(map[int]*msgCodeProcessor),
		}
	})
	return instance
}

func (assistant *assistant) GetProcessor(code int) *msgCodeProcessor {
	return assistant.processorMap[code]
}

func (assistant *assistant) registerProcessor(processor *msgCodeProcessor) {
	assistant.processorMap[processor.Code] = processor
}

type msgCodeProcessor struct {
	Code    int
	Handler Handler
	Message message.IMessage
}

func newMsgProcessor(code int, handler Handler, message message.IMessage) *msgCodeProcessor {
	return &msgCodeProcessor{
		Code:    code,
		Handler: handler,
		Message: message,
	}
}

func PreLoad() {
	// 视频初始化注册
	GetAssistantInstance().registerProcessor(newMsgProcessor(0x1000, NewQQMovieHandler(), &message.MovieMessage{}))
}
