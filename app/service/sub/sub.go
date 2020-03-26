package sub

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)
import "github.com/nats-io/stan.go"

var (
	curCount = 0
	maxCount = 2
	rwMutex  = sync.RWMutex{}
)

// 速率匹配
func Init(clusterID, clientID string) {
	//clusterID := "cluster-crawler-server"
	//clientID := "server1"
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		panic(err)
	}
	subject := "crawler_task"

	// Subscribe with manual ack mode and a max in-flight limit of 2
	sc.Subscribe(subject, func(m *stan.Msg) {
		rwMutex.Lock()
		curCount++
		fmt.Printf("[%s]Received message #: %s,count: %d \n", clientID, string(m.Data), curCount)
		rwMutex.Unlock()

		go processTask(m)

		// Message delivery will suspend when the number of unacknowledged messages reaches 2
	}, stan.SetManualAckMode(), stan.MaxInflight(maxCount))
}

func Init2(clusterID, clientID string) {
	//clusterID := "cluster-crawler-server"
	//clientID := "server1"
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		panic(err)
	}
	subject := "crawler_task"

	// Subscribe with manual ack mode and a max in-flight limit of 25
	sc.Subscribe(subject, func(m *stan.Msg) {
		rwMutex.RLock()
		flag := curCount == maxCount
		rwMutex.RUnlock()
		if flag {
			//不去ack,等到goprocess处理完毕再ack todo 这里实际上是已经接受到消息了，但是没有ack，直到触发ack未回复的超时后重发消息从而保证每个消息都处理了(前提是pub使用同步方式投递消息)；正确限流(速率匹配的)的用法应该看上面Init的用法，而不是自己手动设置变量来控制。
			fmt.Println("wait:", string(m.Data), curCount)
		} else {
			rwMutex.Lock()
			curCount++
			fmt.Printf("Received message #: %s,count: %d \n", string(m.Data), curCount)
			//err := m.Ack()
			//if err != nil {
			//	panic(err)
			//}
			rwMutex.Unlock()
			go processTask(m)
		}

	}, stan.SetManualAckMode())
}

func processTask(msg *stan.Msg) {
	// process
	processHelper(msg.Data)

	//ack and recount when process finish
	rwMutex.Lock()
	curCount--
	fmt.Println("processTask fin:", string(msg.Data), "count:", curCount)
	rwMutex.Unlock()
	err := msg.Ack()
	if err != nil {
		panic(err)
	}

}

//todo 这块处理函数可以暴露给调用者来具体编写爬取业务，可以考虑采用模板方法设计模式。对外暴露接口
func processHelper(bytes []byte) {
	rootUrl := string(bytes)
	//拿url的html，解析总页数，生成其他页码的url存入临时list，接下来就是分别访问这些url，然后保存所需要的搜索结果源url。
	urls := getUrls(rootUrl)

	var wg sync.WaitGroup

	for _, url := range urls {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			parse(url)
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
}

func getUrls(s string) []string {
	//https://www.baidu.com/s?wd=go&pn=30
	prefix := s //todo s 需要是https://www.baidu.com/s?wd=go
	var urls []string
	for i := 0; i < 10; i++ {
		urls = append(urls, prefix+"&pn="+strconv.Itoa(i*10))
	}

	return urls
}

func parse(url string) {
	//todo 保存url的html

	//这里模拟请求即可
	rand.Seed(time.Now().UnixNano())
	s := rand.Intn(5)
	time.Sleep(time.Duration(s) * time.Second)
	fmt.Println(url, s)
	// htmlText := http.Get(url)
	// save(htmlText)
}
