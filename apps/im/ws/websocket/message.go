package websocket

type FrameType uint8

const (
	FrameData FrameType = 0x0
	FramePing FrameType = 0x1
	FrameErr  FrameType = 0x2
)

type Message struct {
	FrameType `json:"frameType"`
	Method    string      `json:"method,omitempty"`
	FormId    string      `json:"formId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// NewErrMessage 创建一个错误消息
func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}

// NewMessage 创建一个数据消息
func NewMessage(fid string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    fid,
		Data:      data,
	}
}
