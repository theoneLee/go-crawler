package factory

import (
	"errors"
	"github.com/nats-io/stan.go"
	"go-crawler/service/sub/handler"
	"go-crawler/service/sub/template_pkg"
	"strings"
)

//todo 根据消息类别创建不同的Consumer ，这里简单的使用字符串判断
func GetConsumer(msg *stan.Msg) (template_pkg.Consumer, error) {
	str := string(msg.Data)
	if strings.Contains(str, "www.baidu.com") {
		return handler.NewBaiduConsumer(), nil
	}
	if strings.Contains(str, "www.google.com") {
		return handler.NewGoogleConsumer(), nil
	}

	return nil, errors.New("error consumer type")
}
