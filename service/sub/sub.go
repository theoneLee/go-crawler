package sub

import (
	"fmt"
	"go-crawler/service/sub/eumn"
	"go-crawler/service/sub/factory"
	"log"
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
