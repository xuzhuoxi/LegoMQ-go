package consumer

import (
	"fmt"
	"testing"
)

func TestPrintConsumer_ConsumeMessage(t *testing.T) {
	consumer, err := PrintConsumer.NewMessageConsumer()
	if nil != err {
		t.Fatal(err)
	}
	err = consumer.ConsumeMessage(msgNil)
	if nil != err {
		fmt.Println("Err1:", err)
	}
	err = consumer.ConsumeMessage(msgEmpty)
	if nil != err {
		fmt.Println("Err2:", err)
	}
	err = consumer.ConsumeMessage(msgDefault)
	if nil != err {
		fmt.Println("Err3:", err)
	}
}
