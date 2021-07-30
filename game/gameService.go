package game

import (
	"fmt"
	"github.com/duanhf2012/origin/node"
	"github.com/duanhf2012/origin/service"
)

func init() {
	node.Setup(&GameService{})
}

type GameService struct {
	service.Service
}

func (slf *GameService) OnInit() error {
	slf.RegRawRpc(1, &LogicHandler{})

	//监听其他Node结点连接和断开事件
	slf.RegRpcListener(slf)
	return nil
}

type InputData struct {
	A int
	B int
}

func (slf *GameService) OnNodeConnected(nodeId int) {
	fmt.Printf("node id %d is conntected.\n", nodeId)
}

func (slf *GameService) OnNodeDisconnect(nodeId int) {
	fmt.Printf("node id %d is disconntected.\n", nodeId)
}

func (slf *GameService) RPC_Sum(input *InputData, output *int) error {
	*output = input.A + input.B
	return nil
}
