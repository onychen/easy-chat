package websocket

type FrameType uint8

const (
	FrameData FrameType = 0x0
	FramePing FrameType = 0x1
)

type Message struct {
	FrameType `json:"frameType"`
	Method    string      `json:"method,omitempty"`
	UserId    string      `json:"userId,omitempty"`
	FormId    string      `json:"formId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func NewMessage(fid string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    fid,
		Data:      data,
	}
}
