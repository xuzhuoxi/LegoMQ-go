package consumer

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/logx"
	"testing"
)

func TestLogConsumer_ConsumeMessage(t *testing.T) {
	consumer, err := LogConsumer.NewMessageConsumer()
	if nil != err {
		t.Fatal(err)
	}
	logConsumer := consumer.(ILogMessageConsumer)
	logConsumer.SetConsumerLevel(logx.LevelInfo)
	err = consumer.ConsumeMessage(msgNil)
	fmt.Println("Err1:", err)
	err = consumer.ConsumeMessage(msgEmpty)
	fmt.Println("Err2:", err)
	err = consumer.ConsumeMessage(msgDefault)
	fmt.Println("Err3:", err)
}
