package handler

import (
	"fmt"
	"go-crawler/service/sub/template_pkg"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//todo 新增一个新的爬取网址，
// 只需要在factory构造一个的XXXConsumer和
// 新建一个go文件定义一个新的XXXConsumer，实现Implement接口（即Process(bytes []byte)

type BaiduConsumer struct {
	*template_pkg.Template
}

func NewBaiduConsumer() template_pkg.Consumer {
	c := &BaiduConsumer{}
	template := template_pkg.NewTemplate(c)
	c.Template = template
	return c
}

func (d *BaiduConsumer) Process(bytes []byte) {
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
