package processor

import (
	"encoding/binary"
	"fmt"
	"github.com/duanhf2012/origin/network"
	"github.com/gogo/protobuf/proto"
	"reflect"
)

type MessageInfo struct {
	msgType    reflect.Type
	msgHandler MessageHandler
}

type MessageHandler func(clientId uint64, msg proto.Message, handlerId uint32)
type ConnectHandler func(clientId uint64)
type UnknownMessageHandler func(clientId uint64, msg []byte)

const MsgHeadSize = 8

type PBProcessor struct {
	mapMsg       map[uint16]MessageInfo
	LittleEndian bool

	unknownMessageHandler UnknownMessageHandler
	connectHandler        ConnectHandler
	disconnectHandler     ConnectHandler
	network.INetMempool
}

type PBPackInfo struct {
	//DataLen uint32 //消息的长度
	Id     uint32 //协议的ID (指明协议结构体)
	Mid    uint32 //客户端给的id
	RawMsg []byte //消息的内容
	Msg    proto.Message
}

func NewPBProcessor() *PBProcessor {
	processor := &PBProcessor{mapMsg: map[uint16]MessageInfo{}}
	processor.INetMempool = network.NewMemAreaPool()
	return processor
}

func (pbProcessor *PBProcessor) SetByteOrder(littleEndian bool) {
	pbProcessor.LittleEndian = littleEndian
}

func (slf *PBPackInfo) GetPackType() uint16 {
	return uint16(slf.Id)
}

func (slf *PBPackInfo) GetMsg() proto.Message {
	return slf.Msg
}

// must goroutine safe
func (pbProcessor *PBProcessor) MsgRoute(msg interface{}, userdata interface{}) error {
	pPackInfo := msg.(*PBPackInfo)
	v, ok := pbProcessor.mapMsg[uint16(pPackInfo.Id)]
	if ok == false {
		return fmt.Errorf("Cannot find msgtype %d is register!", pPackInfo.Id)
	}

	v.msgHandler(userdata.(uint64), pPackInfo.Msg, pPackInfo.Id)
	return nil
}

// must goroutine safe
func (pbProcessor *PBProcessor) Unmarshal(data []byte) (interface{}, error) {
	defer pbProcessor.ReleaseByteSlice(data)
	pbMsg := &PBPackInfo{}
	if pbProcessor.LittleEndian == true {
		//pbMsg.DataLen = binary.LittleEndian.Uint32(data[:4])
		pbMsg.Id = binary.LittleEndian.Uint32(data[:4])
		pbMsg.Mid = binary.LittleEndian.Uint32(data[4:8])
	} else {
		//pbMsg.DataLen = binary.BigEndian.Uint32(data[:4])
		pbMsg.Id = binary.BigEndian.Uint32(data[:4])
		pbMsg.Mid = binary.BigEndian.Uint32(data[4:8])
	}
	info, ok := pbProcessor.mapMsg[uint16(pbMsg.Id)]
	if ok == false {
		return nil, fmt.Errorf("cannot find register %d msgtype!", pbMsg.Id)
	}
	msg := reflect.New(info.msgType.Elem()).Interface()
	pbMsg.Msg = msg.(proto.Message)
	err := proto.Unmarshal(data[8:], pbMsg.Msg)
	return pbMsg, err
}

// must goroutine safe
func (pbProcessor *PBProcessor) Marshal(msg interface{}) ([]byte, error) {
	pMsg := msg.(*PBPackInfo)

	var err error
	if pMsg.Msg != nil {
		pMsg.RawMsg, err = proto.Marshal(pMsg.Msg)
		if err != nil {
			return nil, err
		}
	}
	buff := make([]byte, MsgHeadSize, len(pMsg.RawMsg)+MsgHeadSize)
	if pbProcessor.LittleEndian == true {
		binary.LittleEndian.PutUint32(buff[:4], pMsg.Id)
		binary.LittleEndian.PutUint32(buff[4:8], pMsg.Mid)
	} else {
		binary.BigEndian.PutUint32(buff[:4], pMsg.Id)
		binary.BigEndian.PutUint32(buff[4:8], pMsg.Mid)
	}

	buff = append(buff, pMsg.RawMsg...)
	return buff, nil
}

func (pbProcessor *PBProcessor) Register(msgtype uint16, msg proto.Message, handle MessageHandler) {
	var info MessageInfo

	info.msgType = reflect.TypeOf(msg.(proto.Message))
	info.msgHandler = handle
	pbProcessor.mapMsg[msgtype] = info
}

func (pbProcessor *PBProcessor) MakeMsg(msgType uint16, protoMsg proto.Message) *PBPackInfo {
	return &PBPackInfo{Id: uint32(msgType), Msg: protoMsg}
}

func (pbProcessor *PBProcessor) MakeRawMsg(msgType uint16, msg []byte) *PBPackInfo {
	return &PBPackInfo{Id: uint32(msgType), RawMsg: msg}
}

func (pbProcessor *PBProcessor) UnknownMsgRoute(msg interface{}, userData interface{}) {
	pbProcessor.unknownMessageHandler(userData.(uint64), msg.([]byte))
}

// connect event
func (pbProcessor *PBProcessor) ConnectedRoute(userData interface{}) {
	pbProcessor.connectHandler(userData.(uint64))
}

func (pbProcessor *PBProcessor) DisConnectedRoute(userData interface{}) {
	pbProcessor.disconnectHandler(userData.(uint64))
}

func (pbProcessor *PBProcessor) RegisterUnknownMsg(unknownMessageHandler UnknownMessageHandler) {
	pbProcessor.unknownMessageHandler = unknownMessageHandler
}

func (pbProcessor *PBProcessor) RegisterConnected(connectHandler ConnectHandler) {
	pbProcessor.connectHandler = connectHandler
}

func (pbProcessor *PBProcessor) RegisterDisConnected(disconnectHandler ConnectHandler) {
	pbProcessor.disconnectHandler = disconnectHandler
}
