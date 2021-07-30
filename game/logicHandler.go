package game

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type BaseReq struct {
	Clientid  uint64
	HandlerId uint32
	Body      interface{}
}

type REQ struct {
	Msg string
}

type LogicHandler struct {
}

func (lh *LogicHandler) Unmarshal(data []byte) (interface{}, error) {
	reqData := BaseReq{
		Clientid:  binary.BigEndian.Uint64(data[:8]),
		HandlerId: binary.BigEndian.Uint32(data[8:12]),
	}
	body := REQ{}
	err := json.Unmarshal(data[12:], &body)
	reqData.Body = body
	return reqData, err
}

func (lh *LogicHandler) CB(data interface{}) {

	fmt.Println("1101010=>>>>>", data)
}
