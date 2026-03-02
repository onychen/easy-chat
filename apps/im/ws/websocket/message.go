package websocket

import (
	"time"
)

type FrameType uint8

const (
	FrameData  FrameType = 0x0
	FramePing  FrameType = 0x1
	FrameAck   FrameType = 0x2
	FrameNoAck FrameType = 0x3
	FrameCAck  FrameType = 0x4
	FrameErr   FrameType = 0x9
)

type Message struct {
	FrameType `json:"frameType"`
	Id        string      `json:"id"`
	AckSeq    int         `json:"ackSeq,omitempty"`
	ackTime   time.Time   `json:"-"`
	errCount  int         `json:"-"`
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
