package main

import (
	"go-crawler/app/service/pub"
	"go-crawler/app/service/sub"
	"time"
)

func main() {

	clusterID := "test-cluster" //"cluster-crawler-server" todo panic: stan: connect request timeout (possibly wrong cluster ID?)报错是因为nats stream启动时使用的是默认的clusterID 是"test-cluster"

	//go sub.Init2(clusterID,"server1")

	go sub.Init(clusterID, "server1")
	go sub.Init(clusterID, "server2") //允许有任意多个处理节点（sub）

	time.Sleep(100 * time.Millisecond)
	go pub.Init(clusterID, "client")

	select {}
}
