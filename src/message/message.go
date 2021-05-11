package message

type IMessage interface {
	DoDecoder(body []byte) error
	GetMsgCode() int
	GetDownField() string
	GetTranslateField() string
	GetSegmentField() string
}
