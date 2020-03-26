package pub

import (
	"github.com/nats-io/stan.go"
	"strconv"
)

func Init(clusterID, clientID string) {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		panic(err)
	}

	subject := "crawler_task"

	// todo 每隔3s 投递一堆数据
	list := getList()
	for _, value := range list {
		sc.Publish(subject, []byte(value)) // does not return until an ack has been received from NATS Streaming

		//ackHandler := func(ackedNuid string, err error) {
		//	if err != nil {
		//		log.Printf("Warning: error publishing msg id %s: %v\n", ackedNuid, err.Error())
		//	} else {
		//		log.Printf("Received ack for msg id %s\n", ackedNuid)
		//	}
		//}
		//nuid, err := sc.PublishAsync(subject, []byte(value), ackHandler) // returns immediately
		//if err != nil {
		//	log.Printf("Error publishing msg %s: %v\n", nuid, err.Error())
		//}
	}
}

func getList() []string {
	prefix := "https://www.baidu.com/s?wd=go"
	var urls []string
	for i := 0; i < 10; i++ {
		urls = append(urls, prefix+strconv.Itoa(i))
	}

	return urls
}
