package connector

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/duanhf2012/origin/rpc"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	msgpb "origin-jq-server/common/proto/msg"
	"origin-jq-server/processor"
	"strconv"
)

func (slf *ConnService) DpRegister() {
	//注册监听消息类型MsgType_MsgReq，并注册回调
	slf.processor.Register(uint16(msgpb.MsgType_MsgReq), &msgpb.Req{}, slf.OnRequest2)
}

func (slf *ConnService) OnRequest(clientid uint64, msg proto.Message, handlerId uint32) {
	//解析客户端发过来的数据
	pReq := msg.(*msgpb.Req)
	fmt.Println("收到协议!", pReq)
	input := struct {
		A int
		B int
	}{
		A: 300,
		B: 600,
	}
	var output int
	_ = slf.Call("GameService.RPC_Sum", &input, &output)

	_ = slf.tcpService.SendMsg(clientid, &processor.PBPackInfo{
		Id: uint32(msgpb.MsgType_MsgReq),
		Msg: &msgpb.Req{
			Msg: "123456789==>>" + strconv.Itoa(output),
		},
	})
	fmt.Printf("AsyncCall output %d\n", output)

}

type RawInputArgs struct {
	rawData       []byte
	additionParam []byte
}

func (args RawInputArgs) DoFree() {
}

func (args RawInputArgs) DoEscape() {

}

func (args RawInputArgs) GetRawData() []byte {
	return args.rawData
}

func (slf *ConnService) OnRequest2(clientid uint64, msg proto.Message, handlerId uint32) {
	var inputArgs RawInputArgs
	pbMarshaler := jsonpb.Marshaler{}
	_buffer := new(bytes.Buffer)
	err := pbMarshaler.Marshal(_buffer, msg)
	if err != nil {
		return
	}
	data := _buffer.Bytes()
	tempData := make([]byte, len(data)+12)
	binary.BigEndian.PutUint64(tempData[:8], clientid)
	binary.BigEndian.PutUint32(tempData[8:12], handlerId)
	copy(tempData[12:], data)
	inputArgs.rawData = tempData
	_ = slf.RawGoNode(rpc.RpcProcessorGoGoPB, 3, 1, "GameService", inputArgs)
}
