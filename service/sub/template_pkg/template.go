package template_pkg

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"go-crawler/service/sub/eumn"
)

/*
模版方法模式使用继承机制，把通用步骤和通用方法放到父类中，把具体实现延迟到子类中实现。使得实现符合开闭原则。

因为Golang不提供继承机制，需要使用匿名组合模拟实现继承。

此处需要注意：因为父类需要调用子类方法，所以子类需要匿名组合父类的同时，父类需要持有子类的引用。
*/

type Consumer interface {
	Consume(msg *stan.Msg) //通用步骤
}

type Template struct {
	Implement
	//uri string
}

type Implement interface {
	Process(bytes []byte)
}

func (t *Template) Process(bytes []byte) {
	fmt.Println("default process")
}

func NewTemplate(impl Implement) *Template {
	return &Template{
		Implement: impl,
	}
}

func (t *Template) Consume(msg *stan.Msg) {
	//t.uri = uri
	t.Implement.Process(msg.Data) //由子类实现的特殊步骤

	//ack and recount when process finish
	eumn.RwMutex.Lock()
	eumn.CurCount--
	fmt.Println("processTask fin:", string(msg.Data), "count:", eumn.CurCount)
	eumn.RwMutex.Unlock()
	err := msg.Ack()
	if err != nil {
		panic(err)
	}
}
