package connector

import (
	"fmt"
	"github.com/duanhf2012/origin/node"
	"github.com/duanhf2012/origin/service"
	"github.com/duanhf2012/origin/sysservice/tcpservice"
	"origin-jq-server/processor"
)

func init() {
	//因为与gateway中使用的TcpService不允许重复,所以这里使用自定义服务名称
	tcpService := &tcpservice.TcpService{}
	//tcpService.SetName("MyTcpService")
	node.Setup(tcpService)
	node.Setup(&ConnService{})
}

//新建自定义服务TestService1
type ConnService struct {
	service.Service
	processor  *processor.PBProcessor
	tcpService *tcpservice.TcpService
}

func (slf *ConnService) OnInit() error {
	//获取安装好了的TcpService对象
	slf.tcpService = node.GetService("TcpService").(*tcpservice.TcpService)
	//新建内置的protobuf处理器，您也可以自定义路由器，比如json，后续会补充
	slf.processor = processor.NewPBProcessor()

	//注册监听客户连接断开事件
	slf.processor.RegisterDisConnected(slf.OnDisconnected)
	//注册监听客户连接事件
	slf.processor.RegisterConnected(slf.OnConnected)
	slf.DpRegister()
	//将protobuf消息处理器设置到TcpService服务中
	slf.tcpService.SetProcessor(slf.processor, slf.GetEventHandler())

	return nil
}

func (slf *ConnService) OnConnected(clientid uint64) {
	fmt.Printf("client id %d connected\n", clientid)
}

func (slf *ConnService) OnDisconnected(clientid uint64) {
	fmt.Printf("client id %d disconnected\n", clientid)
}
