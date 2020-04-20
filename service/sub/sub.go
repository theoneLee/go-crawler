package sub

import (
	"fmt"
	"go-crawler/service/sub/eumn"
	"go-crawler/service/sub/factory"
	"log"
	"time"
)
import "github.com/nats-io/stan.go"

// 速率匹配 模板方法
// todo 之前未重构代码见git log
func Init3(clusterID, clientID string) {
	//clusterID := "cluster-crawler-server"
	//clientID := "server1"
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		panic(err)
	}
	subject := "crawler_task"

	// Subscribe with manual ack mode and a max in-flight limit of 2
	sc.Subscribe(subject, func(m *stan.Msg) {
		eumn.RwMutex.Lock()
		eumn.CurCount++
		fmt.Printf("[%s]Received message #: %s,count: %d \n", clientID, string(m.Data), eumn.CurCount)
		eumn.RwMutex.Unlock()

		//todo 根据消息类别创建不同的Consumer
		consumer, err := factory.GetConsumer(m)
		if err == nil {
			go consumer.Consume(m)
		} else {
			log.Println(err)
			m.Ack()
		}

		// Message delivery will suspend when the number of unacknowledged messages reaches 2
	}, stan.SetManualAckMode(), stan.MaxInflight(eumn.MaxCount))
}

// 速率匹配 模板方法 队列组（每个消息将仅传递给每个队列组的一个订阅者）
func Init4(clusterID, clientID string) {
	//clusterID := "cluster-crawler-server"
	//clientID := "server1"
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		panic(err)
	}
	subject := "crawler_task"
	queueGroup := "consumeGroup"

	fun := func(m *stan.Msg) {
		eumn.RwMutex.Lock()
		eumn.CurCount++
		fmt.Printf("[%s]Received message #: %s,count: %d \n", clientID, string(m.Data), eumn.CurCount)
		eumn.RwMutex.Unlock()

		//todo 根据消息类别创建不同的Consumer
		consumer, err := factory.GetConsumer(m)
		if err == nil {
			go consumer.Consume(m)
		} else {
			log.Println(err)
			//m.Ack() //如果自己处理不了该url，要么就是真的没写对应的处理函数，要么就是旧节点，如果是旧节点，就优雅退出。
			m.Sub.Unsubscribe()
			sc.Close()
			time.Sleep(10 * time.Second)
		}

		// Message delivery will suspend when the number of unacknowledged messages reaches 2
	}

	// Subscribe with manual ack mode and a max in-flight limit of 2
	sc.QueueSubscribe(subject, queueGroup, fun, stan.SetManualAckMode(), stan.MaxInflight(eumn.MaxCount))
}
