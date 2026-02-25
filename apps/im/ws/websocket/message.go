package websocket

type Message struct {
	Method string      `json:"method,omitempty"`
	UserId string      `json:"userId,omitempty"`
	FormId string      `json:"formId,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func NewMessage(fid string, data interface{}) *Message {
	return &Message{
		FormId: fid,
		Data:   data,
	}
}
