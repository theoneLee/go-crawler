package template_pkg

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"go-crawler/service/sub/eumn"
)

type Consumer interface {
	Consume(msg *stan.Msg)
}

type Template struct {
	Implement
	uri string
}

type Implement interface {
	Process(bytes []byte)
}

func NewTemplate(impl Implement) *Template {
	return &Template{
		Implement: impl,
	}
}

func (t *Template) Consume(msg *stan.Msg) {
	//t.uri = uri
	t.Implement.Process(msg.Data)

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
