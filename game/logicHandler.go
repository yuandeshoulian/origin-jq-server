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

type InData struct {
	Clientid  uint64
	HandlerId uint32
	Data      []byte
}

type LogicHandler struct {
	service *GameService
}

func (lh *LogicHandler) Unmarshal(data []byte) (interface{}, error) {
	reqData := BaseReq{
		Clientid:  binary.BigEndian.Uint64(data[:8]),
		HandlerId: binary.BigEndian.Uint32(data[8:12]),
	}
	body := REQ{} //TODO:构造器
	err := json.Unmarshal(data[12:], &body)
	reqData.Body = body
	return reqData, err
}

func (lh *LogicHandler) CB(data interface{}) {
	data2 := data.(BaseReq)
	body := data2.Body.(REQ) //TODO:构造器
	rBody, _ := json.Marshal(body)
	rData := &InData{
		Clientid:  data2.Clientid,
		HandlerId: data2.HandlerId,
		Data:      rBody,
	}
	_ = lh.service.Go("ConnService.RPC_Notify", rData)
	fmt.Println("1101010=>>>>>", data)
	//TODO:根据请求id分发给各个模块
}
