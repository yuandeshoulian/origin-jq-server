package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	msgpb "origin-jq-server/common/proto/msg"
	pbprocessor "origin-jq-server/processor"
	"time"
)

func A() {
	processor := pbprocessor.NewPBProcessor()
	processor.Register(uint16(msgpb.MsgType_MsgReq), &msgpb.Req{}, nil)
	conn, err := net.Dial("tcp", "127.0.0.1:9930")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for {
		starTime := time.Now().Unix()
		fmt.Println("连接成功...........")
		pbMsg := &pbprocessor.PBPackInfo{
			Id: uint32(msgpb.MsgType_MsgReq),
			Msg: &msgpb.Req{
				Msg: "123456789",
			},
		}
		data, err := processor.Marshal(pbMsg)
		tempData := make([]byte, len(data)+2)
		binary.BigEndian.PutUint16(tempData[:2], uint16(len(data)))
		copy(tempData[2:], data)
		_, err = conn.Write(tempData)
		if err != nil {
			fmt.Println("client write err: ", err)
			return
		}

		dataLenByte := make([]byte, 2)
		_, err = io.ReadFull(conn, dataLenByte)
		datalen := binary.BigEndian.Uint16(dataLenByte)
		dataByte := make([]byte, datalen)
		_, err = io.ReadFull(conn, dataByte)
		rep, _ := processor.Unmarshal(dataByte)
		fmt.Println("client back: ", rep.(*pbprocessor.PBPackInfo).Msg.(*msgpb.Req).Msg)
		useTime := time.Now().Unix() - starTime
		if useTime > 2 {
			fmt.Println("超过2秒=================>", useTime)
		}
		fmt.Println(useTime, "结束", err)
		time.Sleep(time.Microsecond * 100)
	}

}

func main() {

	for i := 0; i < 5000; i++ {
		//fmt.Println("========")
		go A()
	}

	select {}

	//	time.Sleep(time.Second*5)
	//}
}
