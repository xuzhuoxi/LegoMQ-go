package consumer

import (
	"fmt"
	"github.com/xuzhuoxi/infra-go/logx"
	"testing"
)

var settings = []ConsumerSetting{
	ConsumerSetting{Id: "1", Mode: ClearConsumer},
	ConsumerSetting{Id: "2", Mode: PrintConsumer},
	ConsumerSetting{Id: "3", Mode: LogConsumer}}

func TestConsumerGroup(t *testing.T) {
	config, group := NewMessageConsumerGroup()
	config.InitConsumerGroup(settings)

	logConsumer, err := group.GetConsumer("3")
	if nil != err {
		t.Fatal(err)
	}
	logConsumer.(ILogMessageConsumer).SetConsumerLevel(logx.LevelInfo)
	logConsumer.(ILogMessageConsumer).GetLogger().SetConfig(logx.LogConfig{Type: logx.TypeConsole, Level: logx.LevelAll})

	err = group.ConsumeMessage(msgNil, "1")
	if nil != err {
		fmt.Println("Err1:", err)
	}
	err = group.ConsumeMessage(msgEmpty, "2")
	if nil != err {
		fmt.Println("Err2:", err)
	}
	err = group.ConsumeMessage(msgDefault, "3")
	if nil != err {
		fmt.Println("Err3:", err)
	}
	err = group.ConsumeMessage(msgDefault, "4")
	if nil != err {
		fmt.Println("Err4:", err)
	}
}
