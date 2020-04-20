package main

import (
	"go-crawler/service/pub"
	"go-crawler/service/sub"
	"time"
)

func main() {

	clusterID := "test-cluster" //"cluster-crawler-server" todo panic: stan: connect request timeout (possibly wrong cluster ID?)报错是因为nats stream启动时使用的是默认的clusterID 是"test-cluster"

	//fanOut(clusterID)
	queueGroup(clusterID)

	time.Sleep(100 * time.Millisecond)
	go pub.Init(clusterID, "client")

	select {}
}

func fanOut(clusterID string) { //fan out，一个消息会在多个处理节点都消费一次
	go sub.Init3(clusterID, "server1")
	go sub.Init3(clusterID, "server2") //允许有任意多个处理节点（sub）
}

func queueGroup(clusterID string) {
	//	https://github.com/nats-io/stan.go#queue-groups
	go sub.Init4(clusterID, "server1") //clientID由运行时environment传入
	go sub.Init4(clusterID, "server2")
	//允许有任意多个处理节点（sub）,实际上在分布式环境中这些函数应该是独立一个应用来运行的.因此可能存在旧节点和新节点的情况

}
