package consumer

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/logx"
	"testing"
)

func TestLogConsumer_ConsumeMessage(t *testing.T) {
	logConsumer := NewConsoleLogConsumer()
	logConsumer.SetConsumerLevel(logx.LevelInfo)

	err := logConsumer.ConsumeMessage(msgNil)
	if nil != err {
		fmt.Println("Err1:", err)
	}
	err = logConsumer.ConsumeMessage(msgEmpty)
	if nil != err {
		fmt.Println("Err2:", err)
	}
	err = logConsumer.ConsumeMessage(msgDefault)
	if nil != err {
		fmt.Println("Err3:", err)
	}
}
