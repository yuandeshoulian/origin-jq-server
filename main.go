package main

import (
	"github.com/duanhf2012/origin/node"
	"time"

	_ "origin-jq-server/connector"
	_ "origin-jq-server/game"
)

func main() {
	//打开性能分析报告功能，并设置10秒汇报一次
	node.OpenProfilerReport(time.Second * 10)
	node.Start()
}
