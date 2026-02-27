package consumer

import (
	"fmt"
	"testing"

	"github.com/xuzhuoxi/infra-go/logx"
)

var settings = []ConsumerSetting{
	{Id: "1", Mode: ClearConsumer},
	{Id: "2", Mode: PrintConsumer},
	{Id: "3", Mode: LogConsumer, Log: ConsumerSettingLog{Level: logx.LevelTrace,
		Config: []logx.LogConfig{{Type: logx.TypeConsole, Level: logx.LevelAll}}}}}

func TestConsumerGroup(t *testing.T) {
	var err error

	config, group := NewMessageConsumerGroup()
	_, err = config.InitConsumerGroup(settings)
	if nil != err {
		fmt.Println(err)
	}

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
